/**
 AudiobookRepository.swift
 Copyright (c) 2020 Callum Kerr-Edwards
 */

class AudiobookRepository: RealmRepository<AudiobookFile> {
    static var global = AudiobookRepository()
}
