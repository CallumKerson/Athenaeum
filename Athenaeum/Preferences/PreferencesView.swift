/**
 PreferencesView.swift
 Copyright (c) 2020 Callum Kerr-Edwards
 Licensed under the MIT license.
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
        window = NSWindow.createStandardWindow(withTitle: "Preferences",
                                               width: 600,
                                               height: 160)
        window.contentView = NSHostingView(rootView: self)
        window.delegate = prefsWindowDelegate
        window.tabbingMode = .disallowed
        prefsWindowDelegate.windowIsOpen = true
        window.makeKeyAndOrderFront(nil)
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
        if preferences.useImportDirectory {
            return Text(preferences.importPath.deSandboxedPath)
        }
        return Text("")
    }

    class PrefsWindowDelegate: NSObject, NSWindowDelegate {
        var windowIsOpen = false

        func windowWillClose(_: Notification) {
            windowIsOpen = false
        }
    }
}

struct PrefsView_Previews: PreviewProvider {
    static var previews: some View {
        PreferencesView()
    }
}
