/**
 NavigationDetailViewModel.swift
 Copyright (c) 2020 Callum Kerr-Edwards
 */

import Cocoa
import Combine
import Foundation

class NavigationDetailViewModel: ObservableObject {
    var audiobook: Audiobook?
    var fixMatchDialogDisplayed: Bool = false
    var isGoodReadsConfigured: Bool = false
    var isImporting: Bool = false
    let objectWillChange = ObservableObjectPublisher()

    private var didStateChangeCancellable: AnyCancellable?
    let store: Store<GlobalAppState>

    init(id: UUID, store: Store<GlobalAppState>) {
        self.store = store
        self.didStateChangeCancellable = self.store.stateSubject.sink(receiveValue: {
            if let incomingAudiobookLoadable = $0.audiobookState.audiobooks[id] {
                if case let .loaded(incomingAudiobook) = incomingAudiobookLoadable {
                    if self.audiobook != incomingAudiobook {
                        self.audiobook = incomingAudiobook
                        self.objectWillChange.send()
                    }
                }
            }
            let goodReadsConfigured = !$0.preferencesState.goodReadsAPIKey.isBlank
            if self.isGoodReadsConfigured != goodReadsConfigured {
                self.isGoodReadsConfigured = goodReadsConfigured
                self.objectWillChange.send()
            }

            let incomingImportState = ($0.audiobookState.audiobooks.values.filter { $0.isLoading }
                .count != 0)
            if self.isImporting != incomingImportState {
                self.isImporting = incomingImportState
                self.objectWillChange.send()
            }

            let incomingFixMatchDialogDisplayed = $0.audiobookState.fixMatchDialogDisplayed
            if self.fixMatchDialogDisplayed != incomingFixMatchDialogDisplayed {
                self.fixMatchDialogDisplayed = incomingFixMatchDialogDisplayed
                self.objectWillChange.send()
            }
        })
    }

    func showFixMatchDialog() {
        self.store.dispatch(action: AudiobookActions.SetFixMatchDialogVisible(visibility: true))
    }
}
