package app

import (
	"errors"
	"fmt"
	"log"
	"mynav/pkg/core"
	"mynav/pkg/tui"
	"os"
	"time"

	"github.com/awesome-gocui/gocui"
	"github.com/gookit/color"
)

type App struct {
	// api instance
	api *core.Api

	// ui instance (wrapper over gocui)
	ui *tui.Tui

	// views
	header     *Header
	workspaces *Workspaces
	topics     *Topics
	sessions   *Sessions
	preview    *Preview
	info       *WorkspaceInfo

	// worker for processing tasks in FIFO and debouncing
	worker *Worker
}

// worker magic numbers
const (
	// size of the worker queue
	defaultWorkerSize = 100

	// time to debounce for worker
	defaultWorkerDebounce = 200 * time.Millisecond
)

// view styles
const (
	onFrameColor  = gocui.ColorWhite
	offFrameColor = gocui.AttrDim | gocui.ColorWhite
	onTitleColor  = gocui.AttrBold | gocui.ColorGreen
	offTitleColor = gocui.AttrBold | gocui.ColorCyan
)

// text styles
var (
	topicNameColor              = color.New(color.FgYellow, color.Bold)
	workspaceNameColor          = color.New(color.FgBlue, color.Bold)
	timestampColor              = color.New(color.FgDarkGray, color.OpItalic)
	sessionMarkerColor          = color.New(color.FgGreen, color.Bold)
	alternateSessionMarkerColor = color.New(color.Magenta, color.Bold)
)

// global a instance
var a *App

// Inits and starts the app.
func Start() {
	a = newApp()
	a.start()
}

// Inits the app.
func newApp() *App {
	a := &App{}
	return a
}

// Starts the app.
func (a *App) start() {
	// run cli and handle args
	newCli().run()

	// init start refresh queue
	a.worker = newRefreshQueue(200*time.Millisecond, defaultWorkerSize)
	go a.worker.Start()

	// init the app
	a.init()

	// run main loop
	defer a.ui.Close()
	err := a.ui.MainLoop()
	if err != nil {
		if !errors.Is(err, gocui.ErrQuit) {
			log.Panicln(err)
		}
	}
}

// Inits the app (api, tui, views).
func (a *App) init() {
	// define small helper functions
	initApp := func() {
		// initialize UI
		a.initUI()

		// refresh (populate data to the views)
		a.refreshInit()

		// update toast
		available, tag := a.api.UpdateAvailable()
		if available {
			toast(fmt.Sprintf("mynav %s is available", tag), toastWarn)
		}
	}
	close := func() {
		// start closing after 3 seconds, and display the close counter for 6
		a.closeAfter(6, 3*time.Second)
	}

	// init tui
	a.ui = tui.NewTui()

	// init temp ui to ask for initialization and report errors
	a.tempUI()

	// init api
	var err error
	a.api, err = core.NewApi("")
	if err != nil {
		toast(err.Error(), toastError)
		close()
		return
	}

	// if api is initialized then we can initialize the app
	if a.api != nil {
		initApp()
		return
	}

	// get current dir
	curDir, err := os.Getwd()
	if err != nil {
		toast(err.Error(), toastError)
		close()
		return
	}

	// ensure current dir is not home directory
	home, err := os.UserHomeDir()
	if err != nil {
		toast(err.Error(), toastError)
		close()
		return
	}
	if home == curDir {
		toast("mynav cannot be initialized in the home directory, closing...", toastError)
		close()
		return
	}

	// ask to initalize, and handle error cases
	alert(func(b bool) {
		if !b {
			toast("mynav needs a directory to initialize", toastError)
			close()
			return
		}

		// reinit the api in this dir
		a.api, err = core.NewApi(curDir)
		if err != nil {
			toast(err.Error(), toastError)
			close()
			return
		}

		// handle nil just in case (should not be nil again)
		if a.api == nil {
			toast("Could not initialize mynav", toastError)
			close()
			return
		}

		// finally initialize
		initApp()
	}, "No configuration found. Would you like to initialize this directory?")
}

// Inits the UI, views.
func (a *App) initUI() {
	// instantiate views
	hv := newHeader()
	tv := newTopicsView()
	wv := newWorkspcacesView()
	pv := newPreview()
	sv := newSessionsView()
	wiv := newWorkspaceInfo()
	a.header = hv
	a.topics = tv
	a.workspaces = wv
	a.sessions = sv
	a.preview = pv
	a.info = wiv

	// set manager functions that render the views
	a.ui.SetManager(func(t *tui.Tui) error {
		hv.render()
		tv.render()
		wv.render()
		sv.render()
		wiv.render()
		pv.render()
		return nil
	})

	// init the views (configs, actions etc...)
	hv.init()
	tv.init()
	wv.init()
	sv.init()
	wiv.init()
	pv.init()

	// set global key bindings
	a.initGlobalKeys()
}

// Initializes a temporary (incomplete) ui for initialization.
func (a *App) tempUI() {
	// set a manager that runs no renders
	a.ui.SetManager(func(t *tui.Tui) error {
		return nil
	})

	// set only quit keymaps
	quit := func() bool {
		return true
	}
	a.ui.KeyBinding(nil).
		SetWithQuit(gocui.KeyCtrlC, quit, "Quit").
		SetWithQuit('q', quit, "Quit").
		SetWithQuit('q', quit, "Quit")
}

// Focuses a given view by also changing styles.
func (a *App) focusView(view *tui.View) {
	a.ui.FocusView(view)

	// for each "focusable" views
	for _, v := range []*tui.View{
		a.topics.view,
		a.workspaces.view,
		a.sessions.view,
	} {
		if v.Name() == view.Name() {
			v.FrameColor = onFrameColor
			v.TitleColor = onTitleColor
		} else {
			v.FrameColor = offFrameColor
			v.TitleColor = offTitleColor
		}
	}
}

// Applies general styles to view.
func (a *App) styleView(v *tui.View) {
	v.TitleColor = offTitleColor
	v.FrameColor = offFrameColor
	v.FrameRunes = tui.ThinFrame
}

// Wrapper over refresh function that doesnt select anything.
func (a *App) refreshAll() {
	a.refresh(nil, nil, nil)
}

// Refreshes all the views.
// Ensures the refresh or topics is done before workspaces but everything else in async.
// if selectTopic, selectWorkspace are not nil, they will be selected in the views.
// if selectSession is not nil, current session will be shown in preview/info, otherwise current workspace will be shown.
// This is a generalized function to allow for refreshing the entire UI.
func (a *App) refresh(selectTopic *core.Topic, selectWorkspace *core.Workspace, selectSession *core.Session) {
	a.worker.Queue(func() {
		// header in async
		go func() {
			a.header.refresh()
			a.ui.Update(func() {
				a.header.render()
			})
		}()

		// sessions in async
		go func() {
			a.sessions.refresh()
			if selectSession != nil {
				// if selectSession is passed, we select session and refresh preview
				a.sessions.selectSession(selectSession)
				a.sessions.refreshPreview()
			}
			a.ui.Update(func() {
				if selectSession != nil {
					a.sessions.showInfo()
					a.preview.render()
				}
				a.sessions.render()
			})
		}()

		// topics and workspaces, not in a seperate goroutine as we are already in the worker
		a.topics.refresh()
		if selectTopic != nil {
			// select topic if passed
			a.topics.selectTopic(selectTopic)
		}

		// render topics
		a.ui.Update(func() {
			a.topics.render()
		})

		// refresh workspaces after topics
		a.workspaces.refresh()

		// select workspace if passed
		if selectWorkspace != nil {
			a.workspaces.selectWorkspace(selectWorkspace)
		}

		if selectSession == nil {
			// if selectSession was not passed we show the workspace
			a.workspaces.refreshPreview()
		}

		// render workspaces
		a.ui.Update(func() {
			// if selectSession is nil we show workspace (if it was not nil sessions are shown above)
			if selectSession == nil {
				a.workspaces.showInfo()
				a.preview.render()
			}
			a.workspaces.render()
		})
	})
}

// Modified version of refresh designed to run on start up.
// Key difference being setting loading flags, focus initial views and select last workspace
func (a *App) refreshInit() {
	hv := a.header
	tv := a.topics
	wv := a.workspaces
	sv := a.sessions

	a.worker.Queue(func() {
		// header in async
		go func() {
			hv.refresh()
			a.ui.Update(func() {
				hv.render()
			})
		}()

		// sessions in async
		go func() {
			a.ui.Update(func() {
				sv.setLoading(true)
			})
			sv.refresh()
			a.ui.Update(func() {
				sv.setLoading(false)
				sv.render()
			})
		}()

		// topics before worskpaces
		tv.refresh()
		selected := a.api.SelectedWorkspace()
		if selected != nil {
			tv.selectTopic(selected.Topic)
		}
		a.ui.Update(func() {
			tv.render()
		})

		// workspaces
		a.ui.Update(func() {
			wv.setLoading(true)
		})
		wv.refresh()
		a.ui.Update(func() {
			wv.setLoading(false)
		})
		if selected != nil {
			wv.selectWorkspace(selected)
		}
		wv.refreshPreview()
		a.ui.Update(func() {
			// render workspaces
			wv.render()
			wv.showInfo()
			a.preview.render()

			// initial focus
			if selected != nil {
				wv.focus()
			} else {
				tv.focus()
			}
		})
	})
}

// Runs f in between a tui suspend-resume allowing other terminal apps to run.
func (a *App) runAction(f func()) {
	tui.Suspend()
	f()
	tui.Resume()
}

// Closes the app after count seconds and displays a ticker as a toast.
func (a *App) closeAfter(count int, delay time.Duration) {
	time.AfterFunc(delay, func() {
		ticker := time.Tick(time.Second)
		for range ticker {
			if count == 0 {
				a.ui.Close()
				os.Exit(0)
			}
			a.ui.Update(func() {
				toast(fmt.Sprintf("Closing in %d seconds...", count), toastWarn)
			})
			count--
		}
	})
}

// Inits the global actions.
func (a *App) initGlobalKeys() {
	quit := func() bool {
		return true
	}
	a.ui.KeyBinding(nil).
		SetWithQuit(gocui.KeyCtrlC, quit, "Quit").
		SetWithQuit('q', quit, "Quit").
		SetWithQuit('q', quit, "Quit").
		Set('?', "Toggle cheatsheet", func() {
		}).
		Set('<', "Cycle preview left", func() {
			a.preview.decrement()
		}).
		Set('>', "Cycle preview right", func() {
			a.preview.increment()
		}).
		Set('s', "Search for a workspace", func() {
			sd := new(*Search[*core.Workspace])
			*sd = search(SearchDialogConfig[*core.Workspace]{
				onSearch: func(s string) ([][]string, []*core.Workspace) {
					// get all workspaces by name containing
					allWorkspaces := a.api.AllWorkspaces()
					workspaces := allWorkspaces.ByNameContaining(s)

					// get all topics with this name containing
					topics := a.api.AllTopics().ByNameContaining(s)

					// collect all workspaces for each of these topics
					for _, t := range topics {
						ws := allWorkspaces.ByTopic(t)
						workspaces = append(workspaces, ws...)
					}

					workspaces = workspaces.RemoveDuplicates().SortedByTopic()
					rows := make([][]string, 0)
					for _, w := range workspaces {
						session := a.api.Session(w)
						sessionStr := "-"
						if session != nil {
							sessionStr = "Yes"
						}
						rows = append(rows, []string{
							w.Name,
							w.Topic.Name,
							sessionStr,
						})
					}

					return rows, workspaces
				},
				onSelect: func(w *core.Workspace) {
					a.topics.selectTopic(w.Topic)
					a.workspaces.refresh()
					a.workspaces.selectWorkspace(w)

					if *sd != nil {
						(*sd).close()
					}

					a.workspaces.focus()
				},
				onSelectDescription: "Go to workspace",
				searchViewTitle:     "Search a workspace",
				tableViewTitle:      "Result",
				tableTitles: []string{
					"Workspace",
					"Topic",
					"Active Session",
				}, tableProportions: []float64{
					0.4,
					0.4,
					0.2,
				},
				colStyles: []color.Style{
					workspaceNameColor,
					topicNameColor,
					sessionMarkerColor,
				},
			})
		})
}
