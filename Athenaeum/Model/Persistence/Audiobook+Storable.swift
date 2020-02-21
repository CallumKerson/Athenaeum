/**
 Audiobook+Storable.swift
 Copyright (c) 2020 Callum Kerr-Edwards
 Licensed under the MIT license.
 */

import Foundation

extension AudiobookFile: Entity {
    private var storableAudiobook: StorableAudiobook {
        let realmAudiobook = StorableAudiobook()
        realmAudiobook.uuid = location.sha256HashOfContents
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
        storableAudiobook
    }
}
