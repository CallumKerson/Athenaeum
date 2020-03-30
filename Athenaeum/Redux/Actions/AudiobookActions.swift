/**
 Actions.swift
 Copyright (c) 2020 Callum Kerr-Edwards
 */

import Foundation

struct AudiobookActions {
    struct StartingImportOfAudiobook: Action {
        let audiobook: AudioBook
    }

    struct SetAudiobook: Action {
        let audiobook: AudioBook
    }

    struct SetSelectedAudiobook: Action {
        let id: UUID?
    }

    struct SetFixMatchDialogVisible: Action {
        let visibility: Bool
    }

    struct ClearErrors: Action {}
}
