/**
 AppDelegate.swift
 Copyright (c) 2020 Callum Kerr-Edwards
 */

import Cocoa
import RealmSwift
import SwiftUI
import SwiftyBeaver

let log = SwiftyBeaver.self

@NSApplicationMain
class AppDelegate: NSObject, NSApplicationDelegate {
    var window: NSWindow!
    var prefsView: PrefsView?

    func applicationDidFinishLaunching(_: Notification) {
        // MARK: Logging

        log.addDestination(ConsoleDestination())

        // MARK: Setting up library

//        FileLibrarian.global.setUpLibraryPath()

        // MARK: Main View

        #if DEBUG
            do {
                try FileManager
                    .default
                    .removeItem(at: Realm.Configuration.defaultConfiguration.fileURL!)
            } catch {
                log.info("Hello")
            }

        #endif

        log.info("Creating main view")
//        let contentView = LibraryView(withRepository: AudiobookRepository
//            .global)

//        let sampleAudiobook = AudioBook(id: UUID(),
//                                        title: "The Way of Kings",
//                                        location: URL(string: "https://www.goodreads.com/book/show/7235533-the-way-of-kings")!)
//        let sampleStore = Store<GlobalAppState>(reducer: appStateReducer,
//                                                middleware: [logMiddleware],
//                                                state: GlobalAppState(audiobookState: AudiobookState(audiobooks: [sampleAudiobook])))

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
        if let prefsView = prefsView,
            prefsView.prefsWindowDelegate.windowIsOpen {
            prefsView.window.makeKeyAndOrderFront(self)
        } else {
            self.prefsView = PrefsView(PrefsViewModel(withStore: store))
        }
    }

    @IBAction func importMenuItemActionHandler(_: NSMenuItem) {
        log.info("Opening import dialog from menu item")
//        FileLibrarian.global.openImportAudiobookDialog()
    }
}
