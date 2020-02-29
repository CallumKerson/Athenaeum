/**
 AudiobookRowViewModel.swift
 Copyright (c) 2020 Callum Kerr-Edwards
 */

import Combine
import Foundation

class AudiobookRowViewModel: ObservableObject {
    var audiobook: AudioBook?
    let objectWillChange = ObservableObjectPublisher()

    private var didStateChangeCancellable: AnyCancellable?
    private var store: Store<GlobalAppState>

    init(id: UUID, store: Store<GlobalAppState>) {
        self.store = store
        self.didStateChangeCancellable = self.store.stateSubject.sink(receiveValue: {
            if let incomingAudiobook = $0.audiobookState.audiobooks[id] {
                if self.audiobook != incomingAudiobook {
                    self.audiobook = incomingAudiobook
                    self.objectWillChange.send()
                }
            }
        })
    }
}
