/**
 NavigationHeaderViewModel.swift
 Copyright (c) 2020 Callum Kerr-Edwards
 */

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
            let incomingImportState = ($0.audiobookState.audiobooks.values.filter { $0.isLoading }
                .count != 0)
            if self.isImporting != incomingImportState {
                self.isImporting = incomingImportState
                self.objectWillChange.send()
            }
        })
    }

    func importButtonAction() {
        importFromOpenDialog(store: self.store)
    }
}
