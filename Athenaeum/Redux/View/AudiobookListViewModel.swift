/**
 AudiobookViewModel.swift
 Copyright (c) 2020 Callum Kerr-Edwards
 */
import Combine
import Foundation

class AudiobookListViewModel: ObservableObject {
    var audiobooks: [AudioBook] = []
    var selectedAudiobook: AudioBook? {
        didSet {
            self.store
                .dispatch(action: AudiobookActions
                    .SetSelectedAudiobook(audiobook: self.selectedAudiobook))
        }
    }

    let objectWillChange = ObservableObjectPublisher()

    private var didStateChangeCancellable: AnyCancellable?
    var store: Store<GlobalAppState>

    init(store: Store<GlobalAppState>) {
        self.store = store
        self.didStateChangeCancellable = self.store.stateSubject.sink(receiveValue: {
            let recievedAudiobooks = Array($0.audiobookState.audiobooks.values)
            if self.audiobooks != recievedAudiobooks {
                self.audiobooks = recievedAudiobooks
                self.objectWillChange.send()
            }
            if let selectedAudiobook = $0.audiobookState.selectedAudiobook {
                if self.selectedAudiobook != selectedAudiobook {
                    self.selectedAudiobook = selectedAudiobook
                    self.objectWillChange.send()
                }
            }
        })
    }
}
