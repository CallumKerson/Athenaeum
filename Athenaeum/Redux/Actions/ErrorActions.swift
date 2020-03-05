/**
 ErrorActions.swift
 Copyright (c) 2020 Callum Kerr-Edwards
 */

import Foundation

struct ErrorActions {
    struct SetImportedFileIsOfWrongTypeError: Action {
        let audiobook: AudioBook
    }

    struct SetImportedFileURLCannotBeOpenedError: Action {
        let audiobook: AudioBook
    }

    struct SetImportedFileAlreadyExistsError: Action {
        let audiobook: AudioBook
    }
}
