/**
 EditViewModel.swift
 Copyright (c) 2023 Callum Kerr-Edwards
 */

import Combine
import Foundation

class EditViewModel: ObservableObject, Identifiable {
    let originalBook: Book
    let booksLogicController: BooksLogicController

    @Published var error: AthenaeumError?
    @Published var processing: Bool = false
    var viewDismissalModePublisher = PassthroughSubject<Bool, Never>()
    private var shouldDismissView = false {
        didSet {
            self.viewDismissalModePublisher.send(self.shouldDismissView)
        }
    }

    var subscriptions: Set<AnyCancellable> = []
    var editSummaryViewModel: EditMultilineTextFieldViewModel
    var editSeriesViewModel: EditSeriesViewModel
    var editGenreViewModel: EditGenreViewModel
    var editTitleViewModel: EditTextFieldViewModel
    var editReleaseDateViewModel: EditDateViewModel

    init(booksLogicController: BooksLogicController, originalBook: Book) {
        self.booksLogicController = booksLogicController
        self.originalBook = originalBook
        self.editTitleViewModel = EditTextFieldViewModel(
            label: "Title",
            initialText: originalBook.title
        )
        self.editReleaseDateViewModel = EditDateViewModel(
            label: "Release Date",
            initialDate: originalBook.releaseDate
        )
        self.editSummaryViewModel = EditMultilineTextFieldViewModel(
            sectionLabel: "Summary",
            initialMultilineText: originalBook.summary
        )
        self.editSeriesViewModel = EditSeriesViewModel(
            initialTitle: originalBook.series?.title,
            initalEntry: originalBook.series?.entry
        )
        self.editGenreViewModel = EditGenreViewModel(initialGenre: originalBook.genre?[0])
    }

    func update() {
        let title = self.editTitleViewModel.text()
        let releaseDate = self.editReleaseDateViewModel.date()
        let summary = self.editSummaryViewModel.multilineText()
        let series = self.editSeriesViewModel.series()
        let genre = self.editGenreViewModel.genre()

        self.processing = true
        let editBook = EditBook(
            title: title,
            summary: summary,
            releaseDate: releaseDate,
            series: series,
            genre: genre
        )
        self.booksLogicController
            .updateBook(withID: self.originalBook.id, updatedBook: editBook)
            .receive(on: DispatchQueue.main)
            .sink(receiveCompletion: { [weak self] value in
                guard let self = self else { return }
                if case let .failure(error) = value {
                    logger.error("Failed to update book details for id \(self.originalBook.id)")
                    self.error = error
                    return
                } else {
                    self.error = nil
                }

                self.shouldDismiss()
                self.processing = false

            }, receiveValue: { [weak self] item in
                guard self != nil else { return }
                logger.info("Updating Book Details for book with id \(item.id)")
            })
            .store(in: &self.subscriptions)
    }

    func shouldDismiss() {
        self.shouldDismissView = true
    }
}
