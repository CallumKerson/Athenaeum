/**
 PrefsView.swift
 Copyright (c) 2020 Callum Kerr-Edwards
 */

import SwiftUI

struct PrefsView: View {
    @State var prefsWindowDelegate = PrefsWindowDelegate()
    @State var libraryPathSelection = 0
    @State var importPathSelection = 0

    @State var autoImport: Bool = false

    @ObservedObject var viewModel: PrefsViewModel

    var window: NSWindow!
    init(_ prefsViewModel: PrefsViewModel) {
        self.viewModel = prefsViewModel
        self.autoImport = prefsViewModel.useAutoImport
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
                        Text(viewModel.libraryPath).tag(0)
            })
            }.padding(.bottom)
            Section {
                Toggle("Use Automatic Import",
                       isOn: $viewModel.useAutoImport)
            }
            Section {
                TextField("GoodReads API Key",
                          text: $viewModel.goodReadsAPIKey)
            }
        }
        .frame(minWidth: 600, maxWidth: 600, minHeight: 160, maxHeight: 160)
        .padding()
    }

    class PrefsWindowDelegate: NSObject, NSWindowDelegate {
        var windowIsOpen = false

        func windowWillClose(_: Notification) {
            self.windowIsOpen = false
        }
    }
}

#if DEBUG
    struct PrefsView_Previews: PreviewProvider {
        static var previews: some View {
            PrefsView(PrefsViewModel(withStore: sampleStore))
        }
    }
#endif
