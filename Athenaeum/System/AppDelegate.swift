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
    var prefsView: PreferencesView<UserDefaultsPreferencesStore>?

    func applicationDidFinishLaunching(_: Notification) {
        // MARK: Logging

        log.addDestination(ConsoleDestination())

        // MARK: Main View

        log.info("Creating main view")

        let contentView = LibraryView(withRepository: AudiobookRepository
            .global)
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
            self
                .prefsView =
                PreferencesView(withPreferences: UserDefaultsPreferencesStore
                        .global)
        }
    }

    @IBAction func importMenuItemActionHandler(_: NSMenuItem) {
        log.info("Opening import dialog from menu item")
        Import(withPreferences: UserDefaultsPreferencesStore.global,
               withRepository: AudiobookRepository.global)
            .openImportAudiobookDialog()
    }
}
