/**
 AppDelegate.swift
 Copyright (c) 2020 Callum Kerr-Edwards
 */

import Cocoa

import SwiftUI
import SwiftyBeaver

let log = SwiftyBeaver.self

@NSApplicationMain
class AppDelegate: NSObject, NSApplicationDelegate {
    var window: NSWindow!
    var store: Store<GlobalAppState>?
    var prefsView: PrefsView?
    var persistence: PersistenceService?
    var metadata: MetadataService?
    var podcast: PodcastFeedService?

    func applicationDidFinishLaunching(_: Notification) {
        // MARK: Logging

        log.addDestination(ConsoleDestination())

        // MARK: Setting app state

        let decoder = JSONDecoder()
        if let jsonData = try? Data(contentsOf: PersistenceService.getSaveURL()),
            let state = try? decoder.decode(GlobalAppState.self, from: jsonData) {
            store = Store<GlobalAppState>(reducer: appStateReducer,
                                          middleware: [logMiddleware],
                                          state: state)
        } else {
            store = Store<GlobalAppState>(reducer: appStateReducer,
                                          middleware: [logMiddleware],
                                          state: GlobalAppState())
        }

        guard let store = store else {
            log.error("Invalid app state store")
            fatalError("Invalid app state store")
        }

        self.persistence = PersistenceService(store: store)
        self.metadata = MetadataService(store: store)
        self.podcast = PodcastFeedService(store: store)

        // MARK: Main View

        log.info("Creating main view")

        let contentView = ContentView(ContentViewModel(store: store))

        // Create the window and set the content view.
        self.window = NSWindow(
            contentRect: NSRect(x: 0, y: 0, width: 480, height: 300),
            styleMask: [.titled, .closable, .miniaturizable, .resizable,
                        .fullSizeContentView],
            backing: .buffered, defer: false
        )
        self.window.center()
        self.window.setFrameAutosaveName("Main Window")
        self.window.contentView = NSHostingView(rootView: contentView)
        self.window.makeKeyAndOrderFront(nil)
    }

    func applicationWillTerminate(_: Notification) {
        log.info("Exiting app")
    }

    @IBAction func preferencesMenuItemActionHandler(_: NSMenuItem) {
        log.info("Opening preferences window from menu item")
        guard let store = store else {
            log.error("Invalid app state store")
            fatalError("Invalid app state store")
        }
        if let prefsView = prefsView,
            prefsView.prefsWindowDelegate.windowIsOpen {
            prefsView.window.makeKeyAndOrderFront(self)
        } else {
            self.prefsView = PrefsView(PrefsViewModel(withStore: store))
        }
    }

    @IBAction func importMenuItemActionHandler(_: NSMenuItem) {
        guard let store = store else {
            log.error("Invalid app state store")
            fatalError("Invalid app state store")
        }
        log.info("Opening import dialog from menu item")
        importFromOpenDialog(store: store)
    }
}
