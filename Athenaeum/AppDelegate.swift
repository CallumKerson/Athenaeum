/**
 AppDelegate.swift
 Copyright (c) 2020 Callum Kerr-Edwards
 Licensed under the MIT license.
 */

import Cocoa
import SwiftUI

@NSApplicationMain
class AppDelegate: NSObject, NSApplicationDelegate {
    var window: NSWindow!
    var prefsView: PreferencesView?

    func applicationDidFinishLaunching(_: Notification) {
        UserDefaults.standard.register(defaults: Preferences.defaultPreferences)
        manageMenus()

        // Create the SwiftUI view that provides the window contents.
        let contentView = LibraryView()

        // Create the window and set the content view.
        window = NSWindow(
            contentRect: NSRect(x: 0, y: 0, width: 480, height: 300),
            styleMask: [.titled, .closable, .miniaturizable, .resizable, .fullSizeContentView],
            backing: .buffered, defer: false
        )
        window.center()
        window.setFrameAutosaveName("Main Window")
        window.contentView = NSHostingView(rootView: contentView)
        window.makeKeyAndOrderFront(nil)

//        let connected = Library.global.connect()
//
//        if !connected {
//            abort()
//        }
    }

    func applicationWillTerminate(_: Notification) {
        // Insert code here to tear down your application
    }

    private final func manageMenus() {
        guard let appMenu = NSApplication.shared.mainMenu?.item(at: 0)?.submenu else { return }
        let i: Int = appMenu.indexOfItem(withTitle: "Preferences…")
        guard let preferencesMenuItem = appMenu.item(at: i) else { return }
        appMenu.removeItem(preferencesMenuItem)
        preferencesMenuItem.action = #selector(preferencesMenuItemActionHandler(_:))
        appMenu.addItem(preferencesMenuItem)

        guard let fileMenu = NSApplication.shared.mainMenu?.item(at: 1)?.submenu else { return }
        let importItem = NSMenuItem()
        importItem.title = "Import"
        importItem.keyEquivalent = "i"
        importItem.keyEquivalentModifierMask = [.command]
        importItem.isEnabled = true
        importItem.action = #selector(importMenuItemActionHandler(_:))

        fileMenu.addItem(importItem)
    }

    @objc private func importMenuItemActionHandler(_: NSMenuItem) {
        openImportAudiobookDialog()
    }

    @objc private func preferencesMenuItemActionHandler(_: NSMenuItem) {
        if let prefsView = prefsView, prefsView.prefsWindowDelegate.windowIsOpen {
            prefsView.window.makeKeyAndOrderFront(self)
        } else {
            prefsView = PreferencesView()
        }
    }
}
