package shutdown

import "sync"

// ShutdownCallback is an interface that you have to implement for callbacks.
// OnShutdown will be called when shutdown is request. The Parameter
// is the name of the ShutdownManager that requested shutdown.
type ShutdownCallback interface {
	OnShudown(string) error
}

//ShutdownFunc is a hepler type, so you can easily provide annoymous function
// as ShutdownCallbacks.
type ShutdownFunc func(string) error

// OnShutdown defines the action needded to run when shutdown triggered.
func (f ShutdownFunc) OnShutdown(ShutdownManager string) error {
	return f(ShutdownManager)
}

// ShutdownManager is an ineterface implemnted by ShutdownManagers.
// Getname return the name of ShutdownManager.
// Start ShutdownMangers start listening for shutdown request.
// When they call StartShutdown on GSInterface,
// firest ShutdownStart() is called, then all ShutdownCallbacks are executed
// and once all ShutdownCallbacks return,ShutdownFinish is called.
type ShutdownManager interface {
	GetName() string
	Start(gs GSInterface) error
	ShutdownStart() error
	ShutdownFinish() error
}

// ErrorHandler is an interface you can pass to SetErrorHandler to
// handle asynchronous errors.
type ErrorHandler interface {
	OnError(err error)
}

// ErrorFunc is a hepler type, so you can easily provide anonymous functions
// as ErrorHandlers.
type ErrorFunc func(err error)

// OnError defines the action need to run when error occurred.
func (f ErrorFunc) OnError(err error) {
	f(err)
}

// GSInterface is an interface implemented by GracefulShutdown,
// that gets passed to ShutdownManger to call StartShutdown when shutdown
// is requested.
type GSInterface interface {
	StartShutdown(sm ShutdownManager)
	ReportError(err error)
	AddShutdownCallback(shutdownCallbak ShutdownCallback)
}

// GraceFulShutdown is main struct that handles ShutdownCallbacks and
// ShutdownManagers. Initialize it with New.
type GraceFulShutdown struct {
	callbacks    []ShutdownCallback
	managers     []ShutdownManager
	errorHandler ErrorHandler
}

// New initializes GraceFulShutdown.
func New() *GraceFulShutdown {
	return &GraceFulShutdown{
		callbacks: make([]ShutdownCallback, 0, 10),
		managers:  make([]ShutdownManager, 0, 3),
	}
}

// Start calls start on all added ShutdownMangers. The Shutdownmanagers
// start to listen to shutdown request.  Returns an error if any ShutdownMangers
// return an error.
func (gs *GraceFulShutdown) Start() error {
	for _, manager := range gs.managers {
		if err := manager.Start(gs); err != nil {
			return err
		}
	}
	return nil
}

// AddShutdownManager adds a ShutdownManger that will listen to shutdown request.
func (gs *GraceFulShutdown) AddShutdownManager(manager ShutdownManager) {
	gs.managers = append(gs.managers, manager)
}

// AddShutdownCallback adds a ShutdownCallback that willbe called
// when shutdown is requested.
//
// AddShutdownCallback(shutdown.ShutdownFunc(func() error {
//		// callback code here.
//		return nil
// }))
func (gs *GraceFulShutdown) AddShutdownCallback(shutdownCallback ShutdownCallback) {
	gs.callbacks = append(gs.callbacks, shutdownCallback)
}

// SetErrorHandler sets an ErrorHandler that will be called when an error
// is encountered in SHutdownCallback or in ShutdownManager.
//
// SetErrorHandler(shutdown.ErrorFunc(func (err error) {
//		// handle error here.
// }))
func (gs *GraceFulShutdown) SetErrorHanler(errorHandler ErrorHandler) {
	gs.errorHandler = errorHandler
}

// StartShutdown is called from a ShutdownManager and will initiate shutdown.
// First call ShutdownStart on ShutdownManger,
// call all ShutdownCallbacks, wait for callbacks to finish and
// call ShutdownFinish on ShutdownManger.
func (gs *GraceFulShutdown) StartShutdown(sm ShutdownManager) {
	gs.ReportError(sm.ShutdownStart())

	var wg sync.WaitGroup

	for _, shutdownCallback := range gs.callbacks {
		wg.Add(1)
		go func(shutdownCallback ShutdownCallback) {
			defer wg.Done()

			gs.ReportError(shutdownCallback.OnShudown(sm.GetName()))
		}(shutdownCallback)
	}

	wg.Wait()

	gs.ReportError(sm.ShutdownFinish())

}

// ReportError is a function that can be used to report errors to
// ErrorHandler. It is used in ShutdownManagers.
func (gs *GraceFulShutdown) ReportError(err error) {
	if err != nil && gs.errorHandler != nil {
		gs.errorHandler.OnError(err)
	}
}
