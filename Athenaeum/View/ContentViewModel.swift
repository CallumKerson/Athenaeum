/**
 ContentViewModel.swift
 Copyright (c) 2020 Callum Kerr-Edwards
 */

import Combine
import Foundation

class ContentViewModel: ObservableObject {
    var selectedAudiobook: Audiobook? {
        didSet {
            self.store
                .dispatch(action: AudiobookActions
                    .SetSelectedAudiobook(id: self.selectedAudiobook?.id))
        }
    }

    let objectWillChange = ObservableObjectPublisher()

    private var didStateChangeCancellable: AnyCancellable?
    let store: Store<GlobalAppState>

    init(store: Store<GlobalAppState>) {
        self.store = store
        self.didStateChangeCancellable = self.store.stateSubject.sink(receiveValue: {
            if let selectedAudiobookID = $0.audiobookState.selectedAudiobookID {
                let selectedAudiobook = $0.audiobookState.audiobooks[selectedAudiobookID]
                if self.selectedAudiobook != selectedAudiobook?.get() {
                    self.selectedAudiobook = selectedAudiobook?.get()
                    self.objectWillChange.send()
                }
            }
        })
    }
}
