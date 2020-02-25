/**
 RealmRepositoryTests.swift
 Copyright (c) 2020 Callum Kerr-Edwards
 */

@testable import Athenaeum
import Combine
import RealmSwift
import XCTest

class RealmRepository🧪Tests: XCTestCase {
    var prideAndPrejudiceURL: URL?
    var theFifthSeasonURL: URL?

    override func setUp() {
        Realm.Configuration.defaultConfiguration.inMemoryIdentifier = name

        let temporaryDirectoryURL = URL(fileURLWithPath: NSTemporaryDirectory(), isDirectory: true)

        prideAndPrejudiceURL = temporaryDirectoryURL.appendingPathComponent("\(UUID().uuidString)PrideAndPrejudice.m4b")
        self.theFifthSeasonURL = temporaryDirectoryURL.appendingPathComponent("\(UUID().uuidString)TheFifthSeason.m4b")

        do {
            try "It is a truth universally acknowledged, that a single man in possession of a good fortune, must be in want of a wife."
                .data(using: .utf8)!.write(to: self.prideAndPrejudiceURL!, options: .atomic)
            try "Let's start with the end of the world, why don't we? Get it over with and move on to more interesting things."
                .data(using: .utf8)!.write(to: self.theFifthSeasonURL!, options: .atomic)
        } catch {
            XCTFail()
        }
    }

    func test_insert_stores_🎧📚_locally() {
        let expectation = XCTestExpectation(description: "Publishes notification of database insert")
        let prideAndPrejudice = AudiobookFile(title: "Pride and Prejudice", author: "Jane Austen", file: prideAndPrejudiceURL!)
        let repository = self.createRepository()

        var 🎧📚: [AudiobookFile] = []

        let cancellable = repository.publisher.sink(receiveValue: { action in
            XCTAssertEqual(action, .insert)
            🎧📚 = repository.getAll()
            expectation.fulfill()
        })

        try! repository.insert(item: prideAndPrejudice)

        XCTAssertNotNil(cancellable)
        wait(for: [expectation], timeout: 1.0)
        XCTAssertEqual(1, 🎧📚.count)
        XCTAssertEqual("Jane Austen", 🎧📚.first?.author)
    }

    func test_update_updated_🎧📚() {
        let expectation = XCTestExpectation(description: "Publishes notification of database update")
        let theFifthSeason = AudiobookFile(title: "Fifth Season", author: "NK Jemisin", file: theFifthSeasonURL!)
        let repository = self.createRepository()
        try! repository.insert(item: theFifthSeason)

        // Proper title and puncutation
        theFifthSeason.title = "The Fifth Season"
        theFifthSeason.author = "N. K. Jemisin"

        var 🎧📚: [AudiobookFile] = []

        let cancellable = repository.publisher.sink(receiveValue: { action in
            XCTAssertEqual(action, .update)
            🎧📚 = repository.getAll()
            expectation.fulfill()
        })

        try! repository.update(item: theFifthSeason)

        XCTAssertNotNil(cancellable)
        XCTAssertEqual("The Fifth Season", 🎧📚.first?.title)
        XCTAssertEqual("N. K. Jemisin", 🎧📚.first?.author)
    }

    func test_delete_removes_🎧📚() {
        let expectation = XCTestExpectation(description: "Publishes notification of database delete")
        let prideAndPrejudice = AudiobookFile(title: "Pride and Prejudice", author: "Jane Austen", file: prideAndPrejudiceURL!)

        let repository = self.createRepository()
        try! repository.insert(item: prideAndPrejudice)

        var 🎧📚: [AudiobookFile] = repository.getAll()
        XCTAssertEqual(1, 🎧📚.count)

        let cancellable = repository.publisher.sink(receiveValue: { action in
            XCTAssertEqual(action, .delete)
            🎧📚 = repository.getAll()
            expectation.fulfill()
        })

        try! repository.delete(item: prideAndPrejudice)

        XCTAssertNotNil(cancellable)
        XCTAssertEqual(0, 🎧📚.count)
    }

    func test_getAll_filters_🎧📚() {
        let theFifthSeason = AudiobookFile(title: "The Fifth Season", author: "N. K. Jemisin", file: theFifthSeasonURL!)
        let prideAndPrejudice = AudiobookFile(title: "Pride and Prejudice", author: "Jane Austen", file: prideAndPrejudiceURL!)

        let repository = self.createRepository()
        try! repository.insert(item: theFifthSeason)
        try! repository.insert(item: prideAndPrejudice)

        let 🎧📚: [AudiobookFile] = repository.getAll(where: NSPredicate(format: "author = %@", theFifthSeason.author))

        XCTAssertEqual(1, 🎧📚.count)
        XCTAssertEqual("The Fifth Season", 🎧📚.first?.title)
    }

    private func createRepository() -> RealmRepository<AudiobookFile> {
        RealmRepository()
    }
}
