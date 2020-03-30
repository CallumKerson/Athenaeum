/**
 PrefsViewModel.swift
 Copyright (c) 2020 Callum Kerr-Edwards
 */

import Combine
import Foundation

class PrefsViewModel: ObservableObject {
    var libraryPath: String = ""
    var useAutoImport: Bool = false {
        didSet {
            self.store
                .dispatch(action: PreferencesActions
                    .UpdateAutoImportPreference(updateValueTo: self.useAutoImport))
            self.objectWillChange.send()
        }
    }

    var goodReadsAPIKey: String = "" {
        didSet {
            store
                .dispatch(action: PreferencesActions
                    .UpdateGoodReadsAPIKeyPreference(updateValueTo: goodReadsAPIKey))
            self.objectWillChange.send()
        }
    }

    let objectWillChange = ObservableObjectPublisher()

    private var didStateChangeCancellable: AnyCancellable?
    private let store: Store<GlobalAppState>

    init(withStore store: Store<GlobalAppState>) {
        self.store = store
        self.didStateChangeCancellable = self.store.stateSubject.sink(receiveValue: {
            if self.libraryPath != $0.preferencesState.libraryURL.deSandboxedPath {
                self.libraryPath = $0.preferencesState.libraryURL.deSandboxedPath
                self.objectWillChange.send()
            }
            if self.useAutoImport != $0.preferencesState.autoImport {
                self.useAutoImport = $0.preferencesState.autoImport
                self.objectWillChange.send()
            }
            if self.goodReadsAPIKey != $0.preferencesState.goodReadsAPIKey {
                self.goodReadsAPIKey = $0.preferencesState.goodReadsAPIKey
                self.objectWillChange.send()
            }
        })
    }
}
