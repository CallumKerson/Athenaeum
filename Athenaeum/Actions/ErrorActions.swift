/**
 ErrorActions.swift
 Copyright (c) 2020 Callum Kerr-Edwards
 */

import Foundation

struct ErrorActions {
    struct SetImportedFileIsOfWrongTypeError: Action {
        let audiobook: Audiobook
    }

    struct SetImportedFileURLCannotBeOpenedError: Action {
        let audiobook: Audiobook
    }

    struct SetImportedFileAlreadyExistsError: Action {
        let audiobook: Audiobook
    }
}
