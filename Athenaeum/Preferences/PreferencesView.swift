/**
 PreferencesView.swift
 Copyright (c) 2020 Callum Kerr-Edwards
 Licensed under the MIT license.
 */

import SwiftUI

struct PreferencesView: View {
    @State var prefsWindowDelegate = PrefsWindowDelegate()
    @State var libraryPathSelection = 0
    @State var useImport: Bool = Preferences.getBool(for: .useImport) ?? false
    @State var importPathSelection = 0
    @State var goodReadsAPIKey = Preferences.getString(for: .goodReadsAPIKey) ?? ""

    var body: some View {
        Form {
            Section {
                Picker(selection: $libraryPathSelection, label:
                    Text("Audiobook Library Path:"),
                       content: {
                        Text(Preferences.libraryPath().deSandboxedPath()).tag(0)
            })
            }.padding(.bottom)

            Section {
                Toggle("Use Import Directory", isOn: Binding(
                    get: {
                        self.useImport
                    },
                    set: { newValue in
                        Preferences.set(newValue, for: .useImport)
                        self.useImport = newValue
                    }
                ))
            }
            Section {
                Picker(selection: $importPathSelection, label:
                    Text("Import Path:"),
                       content: {
                        getImportPath().tag(0)
            })
            }
            .disabled(!useImport)
            .padding(.bottom)
            Section {
                TextField("GoodReads API Key", text: $goodReadsAPIKey)
            }
            .onReceive([self.goodReadsAPIKey].publisher.first()) { newAPIKey in
                log.debug("Updating GoodReads API Key to \(newAPIKey)")
                Preferences.set(newAPIKey, for: .goodReadsAPIKey)
            }
        }
        .frame(minWidth: 600, maxWidth: 600, minHeight: 160, maxHeight: 160)
        .padding()
    }

    func getImportPath() -> Text {
        if useImport {
            return Text(Preferences.importPath().deSandboxedPath())
        }
        return Text("")
    }

    var window: NSWindow!
    init() {
        window = NSWindow.createStandardWindow(withTitle: "Preferences",
                                               width: 600,
                                               height: 160)
        window.contentView = NSHostingView(rootView: self)
        window.delegate = prefsWindowDelegate
        window.tabbingMode = .disallowed
        prefsWindowDelegate.windowIsOpen = true
        window.makeKeyAndOrderFront(nil)
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
