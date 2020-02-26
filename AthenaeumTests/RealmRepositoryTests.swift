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

        let temporaryDirectoryURL = URL(fileURLWithPath: NSTemporaryDirectory(),
                                        isDirectory: true)

        prideAndPrejudiceURL = temporaryDirectoryURL
            .appendingPathComponent("\(UUID().uuidString)PrideAndPrejudice.m4b")
        self.theFifthSeasonURL = temporaryDirectoryURL
            .appendingPathComponent("\(UUID().uuidString)TheFifthSeason.m4b")

        do {
            try "It is a truth universally acknowledged, that a single man in possession of a good fortune, must be in want of a wife."
                .data(using: .utf8)!
                .write(to: self.prideAndPrejudiceURL!, options: .atomic)
            try "Let's start with the end of the world, why don't we? Get it over with and move on to more interesting things."
                .data(using: .utf8)!
                .write(to: self.theFifthSeasonURL!, options: .atomic)
        } catch {
            XCTFail()
        }
    }

    func testInsertItem() {
        /// given
        let expectation =
            XCTestExpectation(description: "Object will change when insert occurs")
        expectation.expectedFulfillmentCount = 2
        let prideAndPrejudice = AudiobookFile(title: "Pride and Prejudice",
                                              author: "Jane Austen",
                                              file: prideAndPrejudiceURL!)
        let repository = self.createRepository()

        let cancellable = repository.objectWillChange.sink { _ in
            expectation.fulfill()
        }

        /// when
        try! repository.insert(item: prideAndPrejudice)

        /// then
        XCTAssertNotNil(cancellable)
        wait(for: [expectation], timeout: 1.0)

        XCTAssertEqual(repository.items.count, 1)
        XCTAssertEqual(repository.items.first?.author, "Jane Austen")
    }

    func testUpdateItem() {
        /// given
        let expectation =
            XCTestExpectation(description: "Object will change when update occurs")
        expectation.expectedFulfillmentCount = 3
        let theFifthSeason = AudiobookFile(title: "Fifth Season",
                                           author: "NK Jemisin",
                                           file: theFifthSeasonURL!)
        let repository = self.createRepository()
        try! repository.insert(item: theFifthSeason)
        let cancellable = repository.objectWillChange.sink { _ in
            expectation.fulfill()
        }

        /// when
        /// Proper title and puncutation
        theFifthSeason.title = "The Fifth Season"
        theFifthSeason.author = "N. K. Jemisin"

        try! repository.update(item: theFifthSeason)

        /// then
        XCTAssertNotNil(cancellable)
        wait(for: [expectation], timeout: 1.0)

        XCTAssertEqual(repository.items.count, 1)
        XCTAssertEqual(repository.items.first?.title, "The Fifth Season")
        XCTAssertEqual(repository.items.first?.author, "N. K. Jemisin")
    }

    func testDeleteItem() {
        /// given
        let insertExpectation =
            XCTestExpectation(description: "Object will change when insert occurs")
        insertExpectation.expectedFulfillmentCount = 2
        let deleteExpectation =
            XCTestExpectation(description: "Object will change when delete occurs")
        deleteExpectation.expectedFulfillmentCount = 3
        let prideAndPrejudice = AudiobookFile(title: "Pride and Prejudice",
                                              author: "Jane Austen",
                                              file: prideAndPrejudiceURL!)

        let repository = self.createRepository()

        let cancellable = repository.objectWillChange.sink { _ in
            insertExpectation.fulfill()
            deleteExpectation.fulfill()
        }

        try! repository.insert(item: prideAndPrejudice)

        XCTAssertNotNil(cancellable)
        wait(for: [insertExpectation], timeout: 1.0)
        XCTAssertEqual(repository.items.count, 1)

        /// when

        try! repository.delete(item: prideAndPrejudice)

        /// then
        XCTAssertNotNil(cancellable)
        wait(for: [deleteExpectation], timeout: 1.0)
        XCTAssertEqual(repository.items.count, 0)
    }

    private func createRepository() -> RealmRepository<AudiobookFile> {
        RealmRepository()
    }
}
