/**
 Audiobook+Storable.swift
 Copyright (c) 2020 Callum Kerr-Edwards
 Licensed under the MIT license.
 */

import Foundation

extension Audiobook: Entity {
    private var storableAudiobook: StorableAudiobook {
        let realmAudiobook = StorableAudiobook()
        realmAudiobook.uuid = file.sha256HashOfContents
        realmAudiobook.title = title
        realmAudiobook.author = author
        realmAudiobook.filePath = file.path
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
