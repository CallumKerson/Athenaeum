/**
 NavigationHeaderViewModel.swift
 Copyright (c) 2020 Callum Kerr-Edwards
 */

import Cocoa
import Combine
import Foundation

class NavigationHeaderViewModel: ObservableObject {
    var isImporting: Bool = false
    let objectWillChange = ObservableObjectPublisher()

    private var didStateChangeCancellable: AnyCancellable?
    var store: Store<GlobalAppState>

    init(store: Store<GlobalAppState>) {
        self.store = store
        self.didStateChangeCancellable = self.store.stateSubject.sink(receiveValue: {
            if self.isImporting != ($0.audiobookState.importsInProgress.count != 0) {
                self.isImporting = ($0.audiobookState.importsInProgress.count != 0)
                self.objectWillChange.send()
            }
        })
    }

    func importFromOpenDialog() {
        let openPanel = NSOpenPanel()
        openPanel.canChooseFiles = true
        openPanel.allowsMultipleSelection = true
        openPanel.canChooseDirectories = false
        openPanel.canCreateDirectories = false
        openPanel.title = "Import Audiobooks"
        openPanel.allowedFileTypes = ["m4b"]

        openPanel.begin { response in
            if response == .OK {
                for url in openPanel.urls {
                    self.store
                        .dispatch(action: AudiobookActions
                            .RequestNewAudiobookFromFile(fileURL: url))
                }
            }
            openPanel.close()
        }
    }
}
