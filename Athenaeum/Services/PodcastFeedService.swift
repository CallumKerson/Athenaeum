/**
 PodcastFeedService.swift
 Copyright (c) 2020 Callum Kerr-Edwards
 */

import Combine
import Foundation
import PodcastFeedKit

struct PodcastFeedService {
    private static let persistenceQueue =
        DispatchQueue(label: "com.umbra.Athenaeum.podcastFeedQueue")

    private var didStateChangeCancellable: AnyCancellable?

    init(store: Store<GlobalAppState>) {
        self.didStateChangeCancellable = store.stateSubject.sink(receiveValue: { state in
            if !state.preferencesState.podcastHostURL.isBlank {
                PodcastFeedService.persistenceQueue.async {
                    do {
                        var episodes: [Episode] = []
                        let recievedAudiobooks: [Audiobook] = Array(state.audiobookState.audiobooks
                            .values)
                            .loadedAudiobooks
                        for book in recievedAudiobooks {
                            if let metadata = book.metadata {
                                episodes.append(try Episode(title: metadata.title,
                                                            publicationDate: metadata
                                                                .publicationDate?
                                                                .asDate ?? Date(),
                                                            audioFile: book.location,
                                                            fileServerLocation: "book.location")
                                        .withSubtitle(metadata.authors?.author)
                                        .withLongSummary(metadata.summary)
                                        .withGUID("\(metadata.title)-\(book.id.uuidString)"))
                            }
                        }

                        let podcast = Podcast(title: "Audiobooks",
                                              link: "\(state.preferencesState.podcastHostURL)/feed.rss")
                            .containsExplicitMaterial()
                            .withLanguage(.englishUK)
                            .withAuthor(state.preferencesState.podcastAuthor)
                            .withOwner(name: state.preferencesState.podcastAuthor,
                                       email: state.preferencesState.podcastEmail)
                            .withImage(link: "\(state.preferencesState.podcastHostURL)/artwork.jpg")
                            .withSummary("A collection of audiobooks")
                            .withCategory(.books)
                            .withSubtitle("Like movies for your mind!")
                            .withEpisodes(episodes
                                .sorted(by: { $0.publicationDate < $1.publicationDate }))

                        let outputFeedLocation = state.preferencesState.libraryURL
                            .appendingPathComponent("feed.rss", isDirectory: false)

                        try podcast.getFeed()
                            .write(to: outputFeedLocation, atomically: true, encoding: .utf8)

                    } catch {
                        log.error("Error writing podcast feed")
                        log.error(error)
                    }
                }
            }
        })
    }
}

private func getDate(from date: String) throws -> Date {
    var variableDate: String = date

    if date.countInstances(of: "-") == 0 {
        variableDate = "\(variableDate)-01-01"
    } else if date.countInstances(of: "-") == 1 {
        variableDate = "\(variableDate)-01"
    }

    let dateFormatter = DateFormatter()
    dateFormatter.locale = Locale(identifier: "en_US_POSIX")
    dateFormatter.dateFormat = "yyyy-MM-dd"

    guard let parsedDate = dateFormatter.date(from: variableDate) else {
        throw DateParseError.notAValidDate(date)
    }
    let calendar = Calendar.current
    var components = calendar.dateComponents([.year, .month, .day], from: parsedDate)
    components.hour = 8
    components.timeZone = TimeZone(identifier: "UTC")!

    guard let caculatedDate = Calendar(identifier: Calendar.Identifier.iso8601)
        .date(from: components)
    else {
        throw DateParseError.notAValidDate(date)
    }
    return caculatedDate
}

enum DateParseError: Error {
    case notAValidDate(_ date: String)
}
