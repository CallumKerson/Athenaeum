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

            // MARK: Validate audiobook file

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

                // MARK: Get audiobook metadata

                let asset = AVURLAsset(url: self.fileURL)
                if let title = asset.commonTitle {
                    newAudiobook.metadata = BookMetadata(title: title)
                    newAudiobook.metadata?.authors = asset.artistsAsAuthors
                    if let date = asset.commonCreationDate {
                        newAudiobook.metadata?.publicationDate = try? PublicationDate(from: date)
                    }
                }

                if let state = getGlobalState(state) {
                    if !state.preferencesState.goodReadsAPIKey.isBlank {
                        if let existingMetadata = newAudiobook.metadata {
                            do {
                                newAudiobook
                                    .metadata = try GoodReadsRESTAPI(apiKey: state.preferencesState
                                        .goodReadsAPIKey)
                                    .getBook(title: existingMetadata.title,
                                             author: existingMetadata.authors?.author)
                            } catch {
                                log
                                    .error("Cannot get goodreads metadata for \(newAudiobook.debugDescription)")
                            }
                        }
                    }
                }

                // MARK: Move audiobook

                if let state = getGlobalState(state) {
                    do {
                        try moveAudiobookToLibrary(&newAudiobook,
                                                   libraryURL: state.preferencesState.libraryURL)
                    } catch {
                        DispatchQueue.main.async {
                            dispatch(ErrorActions
                                .SetImportedFileAlreadyExistsError(audiobook: newAudiobook))
                        }
                        return
                    }
                }
                DispatchQueue.main.async {
                    dispatch(SetAudiobook(audiobook: newAudiobook))
                }
            }
            dispatch(StartingImportOfAudiobook(audiobook: newAudiobook))
        }
    }

    struct UpdateAudiobookFromGoodReads: AsyncAction {
        let goodReadsID: String
        let audiobook: AudioBook

        func execute(state: AppState?, dispatch: @escaping DispatchFunction) {
            // MARK: Validate audiobook file

            AudiobookActions.RequestNewAudiobookFromFile.importQueue.async {
                var updatedAudiobook = self.audiobook

                if let state = getGlobalState(state) {
                    if !state.preferencesState.goodReadsAPIKey.isBlank {
                        if let id = Int(self.goodReadsID) {
                            do {
                                updatedAudiobook
                                    .metadata = try GoodReadsRESTAPI(apiKey: state.preferencesState
                                        .goodReadsAPIKey).getBook(goodReadsID: id)
                            } catch {
                                log
                                    .error("Could not fetch metadata from GoodReads by title and author for audiobook \(self.audiobook)")
                            }
                        } else {
                            log.error("Invalid GoodReads ID \(self.goodReadsID)")
                        }
                    }
                }

                // MARK: Move audiobook

                if let state = getGlobalState(state) {
                    do {
                        try moveAudiobookToLibrary(&updatedAudiobook,
                                                   libraryURL: state.preferencesState.libraryURL)
                    } catch {
                        DispatchQueue.main.async {
                            dispatch(ErrorActions
                                .SetImportedFileAlreadyExistsError(audiobook: updatedAudiobook))
                        }
                        return
                    }
                }
                DispatchQueue.main.async {
                    dispatch(SetAudiobook(audiobook: updatedAudiobook))
                }
            }
            dispatch(StartingImportOfAudiobook(audiobook: self.audiobook))
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
