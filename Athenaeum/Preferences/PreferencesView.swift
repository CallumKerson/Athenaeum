/**
 PreferencesView.swift
 Copyright (c) 2020 Callum Kerr-Edwards
 */

import SwiftUI

struct PreferencesView: View {
    @State var prefsWindowDelegate = PrefsWindowDelegate()
    @State var libraryPathSelection = 0
    @State var importPathSelection = 0

    @ObservedObject var preferences: PreferencesStore

    var window: NSWindow!
    init(withPreferences preferences: PreferencesStore = PreferencesStore.global) {
        self.preferences = preferences
        self.window = NSWindow.createStandardWindow(withTitle: "Preferences",
                                                    width: 600,
                                                    height: 160)
        self.window.contentView = NSHostingView(rootView: self)
        self.window.delegate = self.prefsWindowDelegate
        self.window.tabbingMode = .disallowed
        self.prefsWindowDelegate.windowIsOpen = true
        self.window.makeKeyAndOrderFront(nil)
    }

    var body: some View {
        Form {
            Section {
                Picker(selection: $libraryPathSelection, label:
                    Text("Audiobook Library Path:"),
                       content: {
                        Text(preferences.libraryPath.deSandboxedPath).tag(0)
            })
            }.padding(.bottom)
            Section {
                Toggle("Use Import Directory", isOn: $preferences.useImportDirectory)
            }
            Section {
                Picker(selection: $importPathSelection, label:
                    Text("Import Path:"),
                       content: {
                        getImportPath().tag(0)
            })
            }
            .disabled(!preferences.useImportDirectory)
            .padding(.bottom)
            Section {
                TextField("GoodReads API Key", text: $preferences.goodReadsAPIKey)
            }
        }
        .frame(minWidth: 600, maxWidth: 600, minHeight: 160, maxHeight: 160)
        .padding()
    }

    func getImportPath() -> Text {
        if self.preferences.useImportDirectory {
            return Text(self.preferences.importPath.deSandboxedPath)
        }
        return Text("")
    }

    class PrefsWindowDelegate: NSObject, NSWindowDelegate {
        var windowIsOpen = false

        func windowWillClose(_: Notification) {
            self.windowIsOpen = false
        }
    }
}

struct PrefsView_Previews: PreviewProvider {
    static var previews: some View {
        PreferencesView()
    }
}
