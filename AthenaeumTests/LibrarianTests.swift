/**
 LibrarianTests.swift
 Copyright (c) 2020 Callum Kerr-Edwards
 */

@testable import Athenaeum
import XCTest

class LibrarianTests: XCTestCase {
    override func setUp() {
        Realm.Configuration.defaultConfiguration.inMemoryIdentifier = name
    }

    override func tearDown() {
        // Put teardown code here. This method is called after the invocation of each test method in the class.
    }

    func testSetupOfLibraryDirectoryWhenOneDoesNotExist() {
        /// given
        let temporaryLibraryURL = URL(fileURLWithPath: NSTemporaryDirectory(),
                                      isDirectory: true)
            .appendingPathComponent("\(UUID().uuidString)Library")
        let librarian =
            Librarian(withPreferences: MockPreferences(mockLibraryURL: temporaryLibraryURL),
                      withRepository: createRepository())
        XCTAssertFalse(FileManager.default
            .fileExists(atPath: temporaryLibraryURL.path))

        /// when
        librarian.setUpLibraryPath()

        /// then
        var isDir = ObjCBool(true)
        XCTAssertTrue(FileManager.default
            .fileExists(atPath: temporaryLibraryURL.path, isDirectory: &isDir))
    }

    func testSetupOfLibraryDirectoryWhenOneDoesExist() {
        /// given
        let temporaryLibraryURL = URL(fileURLWithPath: NSTemporaryDirectory(),
                                      isDirectory: true)
            .appendingPathComponent("\(UUID().uuidString)Library")
        let librarian =
            Librarian(withPreferences: MockPreferences(mockLibraryURL: temporaryLibraryURL),
                      withRepository: createRepository())
        try! FileManager.default
            .createDirectory(at: temporaryLibraryURL,
                             withIntermediateDirectories: true)
        var isDir = ObjCBool(true)
        XCTAssertTrue(FileManager.default
            .fileExists(atPath: temporaryLibraryURL.path, isDirectory: &isDir))

        /// when
        librarian.setUpLibraryPath()

        /// then
        XCTAssertTrue(FileManager.default
            .fileExists(atPath: temporaryLibraryURL.path, isDirectory: &isDir))
    }

    func testAutoImportOnLibrarySetup() {
        /// given
        let expectation = XCTestExpectation(description: "Updating repository with imported files")
        expectation.expectedFulfillmentCount = 3

        let temporaryLibraryURL = URL(fileURLWithPath: NSTemporaryDirectory(), isDirectory: true)
            .appendingPathComponent("\(UUID().uuidString)Library")
        let preferences = MockPreferences(mockLibraryURL: temporaryLibraryURL)
        let repository = self.createRepository()
        preferences.useImportDirectory = true
        let librarian = Librarian(withPreferences: preferences, withRepository: repository)

        /// create test library directory
        try! FileManager.default
            .createDirectory(at: temporaryLibraryURL, withIntermediateDirectories: true)
        var isDir = ObjCBool(true)
        XCTAssertTrue(FileManager.default
            .fileExists(atPath: temporaryLibraryURL.path, isDirectory: &isDir))

        /// add files to test library directory
        do {
            try "It is a truth universally acknowledged, that a single man in possession of a good fortune, must be in want of a wife."
                .data(using: .utf8)!
                .write(to: temporaryLibraryURL.appendingPathComponent("PrideAndPrejudice.m4b"),
                       options: .atomic)
            try "Let's start with the end of the world, why don't we? Get it over with and move on to more interesting things."
                .data(using: .utf8)!
                .write(to: temporaryLibraryURL.appendingPathComponent("TheFifthSeason.m4b"),
                       options: .atomic)
        } catch {
            XCTFail()
        }

        /// when
        librarian.setUpLibraryPath()

        /// then
        let cancellable = repository.objectWillChange.sink { _ in
            print("----------------______-----________________")
            expectation.fulfill()
        }
        XCTAssertNotNil(cancellable)
        wait(for: [expectation], timeout: 5.0)
        XCTAssertEqual(repository.items.count, 2)
        let titles = repository.items.map { $0.title }
        XCTAssertTrue(titles.contains("PrideAndPrejudice"))
        XCTAssertTrue(titles.contains("TheFifthSeason"))
        XCTAssertTrue(FileManager.default
            .fileExists(atPath: temporaryLibraryURL
                .appendingPathComponent("Unknown", isDirectory: true)
                .appendingPathComponent("PrideAndPrejudice.m4b").path))
        XCTAssertTrue(FileManager.default
            .fileExists(atPath: temporaryLibraryURL
                .appendingPathComponent("Unknown", isDirectory: true)
                .appendingPathComponent("TheFifthSeason.m4b").path))
        XCTAssertFalse(FileManager.default
            .fileExists(atPath: temporaryLibraryURL.appendingPathComponent("PrideAndPrejudice.m4b")
                .path))
        XCTAssertFalse(FileManager.default
            .fileExists(atPath: temporaryLibraryURL.appendingPathComponent("TheFifthSeason.m4b")
                .path))
    }

    func testMonitoring() {
        /// given
        let expectation = XCTestExpectation(description: "Updating repository with imported files")
        expectation.expectedFulfillmentCount = 2

        let temporaryLibraryURL = URL(fileURLWithPath: NSTemporaryDirectory(), isDirectory: true)
            .appendingPathComponent("\(UUID().uuidString)Library")
        let preferences = MockPreferences(mockLibraryURL: temporaryLibraryURL)
        let repository = self.createRepository()
        preferences.useImportDirectory = true
        let librarian = Librarian(withPreferences: preferences, withRepository: repository)

        /// create test library directory
        try! FileManager.default
            .createDirectory(at: temporaryLibraryURL, withIntermediateDirectories: true)
        librarian.setUpLibraryPath()

        /// when
        librarian.setUpMonitor()

        /// add files to test library directory
        do {
            try "It is a truth universally acknowledged, that a single man in possession of a good fortune, must be in want of a wife."
                .data(using: .utf8)!
                .write(to: temporaryLibraryURL.appendingPathComponent("PrideAndPrejudice.m4b"),
                       options: .atomic)
        } catch {
            XCTFail()
        }

        /// then
        let cancellable = repository.objectWillChange.sink { _ in
            expectation.fulfill()
        }
        XCTAssertNotNil(cancellable)
        wait(for: [expectation], timeout: 5.0)
        XCTAssertEqual(repository.items.first?.title, "PrideAndPrejudice")
    }

    private func createRepository() -> RealmRepository<AudiobookFile> {
        RealmRepository()
    }

    class MockPreferences: PreferencesStore {
        var libraryPath: URL
        var useImportDirectory: Bool = false
        var goodReadsAPIKey: String = ""

        init(mockLibraryURL url: URL) {
            self.libraryPath = url
        }
    }
}
