/**
 EditView.swift
 Copyright (c) 2023 Callum Kerr-Edwards
 */

import SwiftUI

struct EditView: View {
    @Environment(\.presentationMode) var presentation
    @ObservedObject var viewModel: EditViewModel
    let onSave: () -> Void
    @State private var showError = false

    init(viewModel: EditViewModel, onSave: @escaping () -> Void) {
        self.viewModel = viewModel
        self.onSave = onSave
    }

    var body: some View {
        VStack {
            HStack {
                Text("Editing")
                    .font(.largeTitle)
                    .padding(.bottom)
                Spacer()
            }
            Form {
                Section(header: Text("Info").font(.headline)) {
                    VStack {
                        EditTextFieldView(viewModel: &viewModel.editTitleViewModel)
                        EditDateView(viewModel: &viewModel.editReleaseDateViewModel)
                        EditSeriesView(viewModel: &viewModel.editSeriesViewModel)
                        EditGenreView(viewModel: &viewModel.editGenreViewModel)
                    }
                }
                EditMultilineTextFieldView(viewModel: &viewModel.editSummaryViewModel)
            }
            exitButtons
                .padding(.vertical)
        }
        .padding()
        .frame(minWidth: 600, minHeight: 200, alignment: .leading)
        .onReceive(viewModel.viewDismissalModePublisher) { shouldDismiss in
            if shouldDismiss {
                logger.info("Dismissing Edit view")
                onSave()
                self.presentation.wrappedValue.dismiss()
            }
        }
    }

    var exitButtons: some View {
        HStack {
            Button("Cancel") {
                viewModel.shouldDismiss()
            }
            Spacer()
            if viewModel.processing {
                ProgressView()
            } else if let error = viewModel.error {
                Label(error.description, systemImage: "exclamationmark.triangle")
            } else {
                Button("Done") {
                    viewModel.update()
                }
            }
        }
    }
}

struct EditView_Previews: PreviewProvider {
    static var previews: some View {
        EditView(viewModel: EditViewModel(
            booksLogicController: MockController(),
            originalBook: MockController.book
        )) {
            print("")
        }
    }
}
