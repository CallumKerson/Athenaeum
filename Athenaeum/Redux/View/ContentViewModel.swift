/**
 ContentViewModel.swift
 Copyright (c) 2020 Callum Kerr-Edwards
 */

import Combine
import Foundation

class ContentViewModel: ObservableObject {
    var selectedAudiobook: AudioBook? {
        didSet {
            self.store
                .dispatch(action: AudiobookActions
                    .SetSelectedAudiobook(audiobook: self.selectedAudiobook))
        }
    }

    let objectWillChange = ObservableObjectPublisher()

    private var didStateChangeCancellable: AnyCancellable?
    let store: Store<GlobalAppState>

    init(store: Store<GlobalAppState>) {
        self.store = store
        self.didStateChangeCancellable = self.store.stateSubject.sink(receiveValue: {
            if let selectedAudiobook = $0.audiobookState.selectedAudiobook {
                if self.selectedAudiobook != selectedAudiobook {
                    self.selectedAudiobook = selectedAudiobook
                    self.objectWillChange.send()
                }
            }
        })
    }
}
