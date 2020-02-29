/**
 ErrorActions.swift
 Copyright (c) 2020 Callum Kerr-Edwards
 */

import Foundation

struct ErrorActions {
    struct SetImportedFileIsOfWrongTypeError: Action {
        let idNotProcessed: UUID
        let fileURL: URL
    }

    struct SetImportedFileURLCannotBeOpenedError: Action {
        let idNotProcessed: UUID
        let fileURL: URL
    }

    struct SetImportedFileAlreadyExistsError: Action {
        let idNotProcessed: UUID
        let importedFileURL: URL
        let existingAudiobook: AudioBook?
    }
}
