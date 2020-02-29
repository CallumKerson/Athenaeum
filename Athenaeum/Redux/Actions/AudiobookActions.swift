/**
 Actions.swift
 Copyright (c) 2020 Callum Kerr-Edwards
 */

import Foundation

struct AudiobookActions {
    struct RequestNewAudiobookFromFile: AsyncAction {
        let fileURL: URL

        func execute(state _: AppState?, dispatch: @escaping DispatchFunction) {
            let newAudiobookID = UUID()
            dispatch(StartingImportOfAudiobook(id: newAudiobookID))
            dispatch(ValidateAudiobookFile(newAuidobookID: newAudiobookID, fileURL: self.fileURL))
        }
    }

    struct StartingImportOfAudiobook: Action {
        let id: UUID
    }

    struct SetAudiobook: Action {
        let audiobook: AudioBook
    }

    struct SetSelectedAudiobook: Action {
        let audiobook: AudioBook?
    }
}
