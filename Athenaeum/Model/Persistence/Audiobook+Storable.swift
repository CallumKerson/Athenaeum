/**
 Audiobook+Storable.swift
 Copyright (c) 2020 Callum Kerr-Edwards
 */

import Foundation

extension AudiobookFile: Entity {
    private var storableAudiobook: StorableAudiobook {
        let realmAudiobook = StorableAudiobook()
        realmAudiobook.uuid = UUID().uuidString
        realmAudiobook.title = title
        realmAudiobook.author = author
        realmAudiobook.filePath = location.path
        realmAudiobook.publicationDate = publicationDate
        realmAudiobook.narrator = narrator
        realmAudiobook.summary = summary
        realmAudiobook.seriesTitle = series?.title
        realmAudiobook.seriesEntry = series?.entry
        return realmAudiobook
    }

    func toStorable() -> StorableAudiobook {
        self.storableAudiobook
    }
}
