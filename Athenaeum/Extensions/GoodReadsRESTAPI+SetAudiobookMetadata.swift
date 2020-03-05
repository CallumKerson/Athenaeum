/**
 GoodReadsRESTAPI+SetAudiobookMetadata.swift
 Copyright (c) 2020 Callum Kerr-Edwards
 */

import Foundation
import GoodReadsKit

extension GoodReadsRESTAPI {
    /// Looks up book metadata and adds all found to the audiobook
    /// - Parameter audiobook: audiobook for which the metadata is updated
    func setAudiobookMetadataFromGoodReads(audiobook: inout AudioBook) {
        do {
            if let title = audiobook.title {
                let fetchedBook = try self.getBook(title: title,
                                                   author: audiobook.getAuthorsString())

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

                audiobook.title = fetchedBook.title
                audiobook.authors = authors
                audiobook.publicationDate = fetchedBook.getDateString()
                audiobook.isbn = fetchedBook.isbn
                audiobook.bookDescription = fetchedBook.bookDescription
                audiobook.series = series
            }
        } catch {
            log.error("Could not fetch metadata from goodreads")
        }
    }
}
