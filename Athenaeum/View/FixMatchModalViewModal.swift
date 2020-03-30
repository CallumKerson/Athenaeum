/**
 FixMatchModalViewModal.swift
 Copyright (c) 2020 Callum Kerr-Edwards
 */

import Combine
import Foundation

class FixMatchModalViewModel: ObservableObject {
    let audiobook: Audiobook

    private var didStateChangeCancellable: AnyCancellable?
    var store: Store<GlobalAppState>

    init(audiobook: Audiobook, store: Store<GlobalAppState>) {
        self.store = store
        self.audiobook = audiobook
    }

    func fixMatchButton(goodReadsID: String) {
        self.store.dispatch(action: AudiobookActions
            .SetFixMatchDialogVisible(visibility: false))
        self.objectWillChange.send()
        self.store.dispatch(action: AudiobookActions
            .UpdateAudiobookFromGoodReads(goodReadsID: goodReadsID, audiobook: self.audiobook))
    }

    func cancelButtonAction() {
        self.store.dispatch(action: AudiobookActions
            .SetFixMatchDialogVisible(visibility: false))
        self.objectWillChange.send()
    }
}
