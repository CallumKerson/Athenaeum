/**
 AudiobookActions+GoodReads.swift
 Copyright (c) 2020 Callum Kerr-Edwards
 */

import Foundation
import GoodReadsKit

extension AudiobookActions {
    struct UpdateAudiobookMetadataFromGoodReads: AsyncAction {
        let audiobookToUpdate: AudioBook

        func execute(state: AppState?, dispatch: @escaping DispatchFunction) {
            DispatchQueue.global(qos: .userInitiated).async {
                var updatedAudiobook = self.audiobookToUpdate
                if let state = state {
                    if state is GlobalAppState {
                        let goodReadsAPIKey = (state as! GlobalAppState).preferencesState
                            .goodReadsAPIKey
                        if !goodReadsAPIKey.isBlank {
                            let goodReadsAPI = GoodReadsRESTAPI(apiKey: goodReadsAPIKey)
                            do {
                                let fetchedBook = try goodReadsAPI
                                    .getBook(title: self.audiobookToUpdate.title,
                                             author: self
                                                 .audiobookToUpdate
                                                 .getAuthorsString())

                                var series: Series?
                                if let seriesTitle = fetchedBook.seriesTitle,
                                    let seriesEntry = fetchedBook.seriesEntry {
                                    series = Series(title: seriesTitle,
                                                    entry: String(seriesEntry))
                                }
                                var authors: [Author] = []
                                for fetchedAuthor in fetchedBook.authors {
                                    var names = fetchedAuthor.components(separatedBy: " ")
                                    if let lastName = names.last {
                                        names.removeLast(1)
                                        authors
                                            .append(Author(firstName: names.joined(separator: " "),
                                                           lastName: lastName))
                                    } else {
                                        authors
                                            .append(Author(firstName: nil, lastName: fetchedAuthor))
                                    }
                                }

                                updatedAudiobook.title = fetchedBook.title
                                updatedAudiobook.authors = authors
                                updatedAudiobook.publicationDate = fetchedBook.getDateString()
                                updatedAudiobook.isbn = fetchedBook.isbn
                                updatedAudiobook.bookDescription = fetchedBook.bookDescription
                                updatedAudiobook.series = series
                            } catch {
                                log.error("Could not fetch metadata from goodreads")
                            }
                        }
                    }
                }

                DispatchQueue.main.async {
                    dispatch(AudiobookActions
                        .SetAudiobook(audiobook: updatedAudiobook))
                }
            }
        }
    }
}
