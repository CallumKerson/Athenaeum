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
                    .SetUpdatedAutoImportPreference(updatedValue: self.useAutoImport))
            self.objectWillChange.send()
        }
    }

    var goodReadsAPIKey: String = "" {
        didSet {
            self.store
                .dispatch(action: PreferencesActions
                    .SetUpdatedGoodReadsAPIKeyPreference(updatedValue: self.goodReadsAPIKey))
            self.objectWillChange.send()
        }
    }

    var podcastAuthor: String = "" {
        didSet {
            self.store
                .dispatch(action: PreferencesActions
                    .SetUpdatedPodcastAuthorPreference(updatedValue: self.podcastAuthor))
            self.objectWillChange.send()
        }
    }

    var podcastEmail: String = "" {
        didSet {
            self.store
                .dispatch(action: PreferencesActions
                    .SetUpdatedPodcastEmailPreference(updatedValue: self.podcastEmail))
            self.objectWillChange.send()
        }
    }

    var podcastHostURL: String = "" {
        didSet {
            self.store
                .dispatch(action: PreferencesActions
                    .SetUpdatedPodcastHostURL(updatedValue: self.podcastHostURL))
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
            if self.podcastAuthor != $0.preferencesState.podcastAuthor {
                self.podcastAuthor = $0.preferencesState.podcastAuthor
                self.objectWillChange.send()
            }
            if self.podcastEmail != $0.preferencesState.podcastEmail {
                self.podcastEmail = $0.preferencesState.podcastEmail
                self.objectWillChange.send()
            }
        })
    }
}
