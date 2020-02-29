/**
 AudiobookActions+AVAsset.swift
 Copyright (c) 2020 Callum Kerr-Edwards
 */

import AVFoundation
import Foundation

extension AudiobookActions {
    struct ValidateAudiobookFile: AsyncAction {
        let newAuidobookID: UUID
        let fileURL: URL

        func execute(state: AppState?, dispatch: @escaping DispatchFunction) {
            DispatchQueue.global(qos: .userInitiated).async {
                if self.fileURL.pathExtension != "m4b" {
                    dispatch(ErrorActions
                        .SetImportedFileIsOfWrongTypeError(idNotProcessed: self.newAuidobookID,
                                                           fileURL: self.fileURL))
                    return
                }
                if !AVURLAsset(url: self.fileURL).isPlayable {
                    dispatch(ErrorActions
                        .SetImportedFileURLCannotBeOpenedError(idNotProcessed: self.newAuidobookID,
                                                               fileURL: self
                                                                   .fileURL))
                    return
                }
                log.warning("Calculating hash for \(self.fileURL)")
                guard let newFileHash = self.fileURL.sha256HashOfContents else {
                    DispatchQueue.main.async {
                        dispatch(ErrorActions
                            .SetImportedFileURLCannotBeOpenedError(idNotProcessed: self
                                .newAuidobookID,
                                                                   fileURL: self.fileURL))
                    }
                    return
                }
                log.warning("Calculated hash for \(self.fileURL)")
                if let state = state {
                    if state is GlobalAppState {
                        let contentsHashes = (state as! GlobalAppState).audiobookState.audiobooks
                            .map { $0.value.contentsHash }

                        if contentsHashes.contains(newFileHash) {
                            let existingAudiobook = (state as! GlobalAppState).audiobookState
                                .audiobooks.first(where: { $0.value.contentsHash == newFileHash })
                            DispatchQueue.main.async {
                                dispatch(ErrorActions
                                    .SetImportedFileAlreadyExistsError(idNotProcessed: self
                                        .newAuidobookID,
                                                                       importedFileURL: self
                                            .fileURL,
                                                                       existingAudiobook: existingAudiobook?
                                            .value))
                            }
                            return
                        }
                    }
                }
                DispatchQueue.main.async {
                    dispatch(UpdateAudiobookMetadataFromAVMetadata(audiobookToUpdate:
                        AudioBook(id: self.newAuidobookID,
                                  location: self.fileURL,
                                  contentsHash: newFileHash,
                                  title: self.fileURL.deletingPathExtension().lastPathComponent)))
                }
            }
        }
    }

    struct UpdateAudiobookMetadataFromAVMetadata: AsyncAction {
        let audiobookToUpdate: AudioBook

        func execute(state: AppState?, dispatch: @escaping DispatchFunction) {
            DispatchQueue.global(qos: .userInitiated).async {
                var updatedAudiobook = self.audiobookToUpdate

                let asset = AVURLAsset(url: self.audiobookToUpdate.location)

                let metadata = asset.commonMetadata

                if let item = AVMetadataItem.metadataItems(from: metadata,
                                                           filteredByIdentifier: .commonIdentifierTitle)
                    .first {
                    if let titleString = item.stringValue {
                        updatedAudiobook.title = titleString.removeIllegalCharacters
                    }
                }

                if let item = AVMetadataItem.metadataItems(from: metadata,
                                                           filteredByIdentifier: .commonIdentifierArtist)
                    .first {
                    if let artistString = item.stringValue {
                        let names = artistString.components(separatedBy: " ")
                        if names.count == 2 {
                            updatedAudiobook
                                .authors = [Author(firstName: names.first, lastName: names.last!)]
                        }
                    }
                }

                if let item = AVMetadataItem.metadataItems(from: metadata,
                                                           filteredByIdentifier: .commonIdentifierCreationDate)
                    .first {
                    updatedAudiobook.publicationDate = item.stringValue
                }

                if let state = state {
                    if state is GlobalAppState {
                        if !(state as! GlobalAppState).preferencesState.goodReadsAPIKey.isBlank {
                            DispatchQueue.main.async {
                                dispatch(UpdateAudiobookMetadataFromGoodReads(audiobookToUpdate: updatedAudiobook))
                            }
                            return
                        }
                    }
                }
                DispatchQueue.main.async {
                    dispatch(SetAudiobook(audiobook: updatedAudiobook))
                }
            }
        }
    }
}
