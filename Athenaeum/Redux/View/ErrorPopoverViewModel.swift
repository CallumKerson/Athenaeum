/**
 ErrorPopoverViewModel.swift
 Copyright (c) 2020 Callum Kerr-Edwards
 */

import Combine
import Foundation

class ErrorPopoverViewModel: ObservableObject {
    var errors: [DisplayError] = []
    let store: Store<GlobalAppState>

    private var didStateChangeCancellable: AnyCancellable?

    init(store: Store<GlobalAppState>) {
        self.store = store
        self.didStateChangeCancellable = self.store.stateSubject.sink(receiveValue: {
            let receivedErrors: [DisplayError] = Array($0.audiobookState.audiobooks.values)
                .errors.map { DisplayError(path: $0.0.location.path, message: $0.1) }
            if self.errors != receivedErrors {
                self.errors = receivedErrors
                self.objectWillChange.send()
            }
        })
    }

    func clearErrors() {
        self.store.dispatch(action: AudiobookActions.ClearErrors())
    }
}

struct DisplayError: Equatable, Identifiable {
    var id: String { self.path }
    let path: String
    let message: String
}
