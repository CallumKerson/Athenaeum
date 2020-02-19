/**
 AppDelegate.swift
 Copyright (c) 2020 Callum Kerr-Edwards
 Licensed under the MIT license.
 */

import Cocoa
import SwiftUI
import SwiftyBeaver

let log = SwiftyBeaver.self

@NSApplicationMain
class AppDelegate: NSObject, NSApplicationDelegate {
    var window: NSWindow!
    var prefsView: PreferencesView?

    func applicationDidFinishLaunching(_: Notification) {
        // MARK: Logging

        log.addDestination(ConsoleDestination())

        // MARK: Menus

        log.info("Setting up menu items")
        manageMenus()

        // MARK: Main View

        log.info("Creating main view")

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
    }

    func applicationWillTerminate(_: Notification) {
        log.info("Exiting app")
    }

    private final func manageMenus() {
        log.debug("Setting up preferences menu item")
        guard let appMenu = NSApplication.shared.mainMenu?.item(at: 0)?.submenu else { return }
        let i: Int = appMenu.indexOfItem(withTitle: "Preferences…")
        guard let preferencesMenuItem = appMenu.item(at: i) else { return }
        appMenu.removeItem(preferencesMenuItem)
        preferencesMenuItem.action = #selector(preferencesMenuItemActionHandler(_:))
        appMenu.addItem(preferencesMenuItem)

        log.debug("Setting up import menu item")
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
        Import().openImportAudiobookDialog()
    }

    @objc private func preferencesMenuItemActionHandler(_: NSMenuItem) {
        if let prefsView = prefsView, prefsView.prefsWindowDelegate.windowIsOpen {
            prefsView.window.makeKeyAndOrderFront(self)
        } else {
            prefsView = PreferencesView()
        }
    }
}
