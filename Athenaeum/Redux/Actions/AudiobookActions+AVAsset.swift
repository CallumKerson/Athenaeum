/**
 AudiobookActions+AVAsset.swift
 Copyright (c) 2020 Callum Kerr-Edwards
 */

import AVFoundation
import Foundation
import GoodReadsKit

extension AudiobookActions {
    struct RequestNewAudiobookFromFile: AsyncAction {
        static let importQueue = DispatchQueue(label: "com.umbra.Athenaeum.importQueue")
        let fileURL: URL

        func execute(state: AppState?, dispatch: @escaping DispatchFunction) {
            var newAudiobook = AudioBook(id: UUID(), location: self.fileURL)

            AudiobookActions.RequestNewAudiobookFromFile.importQueue.async {
                if self.fileURL.pathExtension != "m4b" {
                    dispatch(ErrorActions
                        .SetImportedFileIsOfWrongTypeError(audiobook: newAudiobook))
                    return
                }
                if !AVURLAsset(url: self.fileURL).isPlayable {
                    dispatch(ErrorActions
                        .SetImportedFileURLCannotBeOpenedError(audiobook: newAudiobook))
                    return
                }

                if let state = getGlobalState(state) {
                    let books = Array(state.audiobookState.audiobooks.values).loadedAudiobooks
                    if books.hasAudibookWithSameFileAs(self.fileURL) {
                        DispatchQueue.main.async {
                            dispatch(ErrorActions
                                .SetImportedFileAlreadyExistsError(audiobook: newAudiobook))
                        }
                        return
                    }
                }

                let asset = AVURLAsset(url: self.fileURL)
                newAudiobook.title = asset.commonTitle
                newAudiobook.authors = asset.artistsAsAuthors
                newAudiobook.publicationDate = asset.commonCreationDate

                if let state = getGlobalState(state) {
                    if !state.preferencesState.goodReadsAPIKey.isBlank {
                        GoodReadsRESTAPI(apiKey: state.preferencesState.goodReadsAPIKey)
                            .setAudiobookMetadataFromGoodReads(audiobook: &newAudiobook)
                    }
                }
                DispatchQueue.main.async {
                    dispatch(SetAudiobook(audiobook: newAudiobook))
                }
            }
            dispatch(StartingImportOfAudiobook(audiobook: newAudiobook))
        }
    }
}

func getGlobalState(_ state: AppState?) -> GlobalAppState? {
    if let state = state {
        if state is GlobalAppState {
            return (state as! GlobalAppState)
        }
    }
    return nil
}
