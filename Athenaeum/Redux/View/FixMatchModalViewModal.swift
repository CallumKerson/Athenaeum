/**
 FixMatchModalViewModal.swift
 Copyright (c) 2020 Callum Kerr-Edwards
 */

import Combine
import Foundation

class FixMatchModalViewModel: ObservableObject {
    let audiobook: AudioBook

    private var didStateChangeCancellable: AnyCancellable?
    var store: Store<GlobalAppState>

    init(audiobook: AudioBook, store: Store<GlobalAppState>) {
        self.store = store
        self.audiobook = audiobook
//        self.didStateChangeCancellable = self.store.stateSubject.sink(receiveValue: {
//            let incomingFixMatchDialogDisplayed = $0.audiobookState.fixMatchDialogDisplayed
//            if self.fixMatchDialogDisplayed != incomingFixMatchDialogDisplayed {
//                self.fixMatchDialogDisplayed = incomingFixMatchDialogDisplayed
//                self.objectWillChange.send()
//            }
//        })
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
