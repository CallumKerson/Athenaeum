/**
 Actions.swift
 Copyright (c) 2020 Callum Kerr-Edwards
 */

import Foundation

struct AudiobookActions {
    struct StartingImportOfAudiobook: Action {
        let audiobook: Audiobook
    }

    struct SetAudiobook: Action {
        let audiobook: Audiobook
    }

    struct SetSelectedAudiobook: Action {
        let id: UUID?
    }

    struct SetFixMatchDialogVisible: Action {
        let visibility: Bool
    }

    struct ClearErrors: Action {}
}
