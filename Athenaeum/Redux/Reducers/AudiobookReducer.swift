/**
 AudiobookReducer.swift
 Copyright (c) 2020 Callum Kerr-Edwards
 */

import Foundation

func audiobookStateReducer(state: AudiobookState, action: Action) -> AudiobookState {
    var state = state
    switch action {
    case let action as AudiobookActions.SetAudiobook:
        state.audiobooks[action.audiobook.id] = .loaded(action.audiobook)
    case let action as AudiobookActions.StartingImportOfAudiobook:
        state.audiobooks[action.audiobook.id] = .loading(action.audiobook)
    case let action as ErrorActions.SetImportedFileAlreadyExistsError:
        state
            .audiobooks[action.audiobook.id] = .errored(action.audiobook,
                                                        message: "File already exists")
    case let action as ErrorActions.SetImportedFileIsOfWrongTypeError:
        state
            .audiobooks[action.audiobook.id] = .errored(action.audiobook,
                                                        message: "File is of wrong type")
    case let action as ErrorActions.SetImportedFileURLCannotBeOpenedError:
        state
            .audiobooks[action.audiobook.id] = .errored(action.audiobook,
                                                        message: "File cannot be accessed")
    case let action as AudiobookActions.SetSelectedAudiobook:
        state.selectedAudiobook = action.audiobook
    case let action as AudiobookActions.SetFixMatchDialogVisible:
        state.fixMatchDialogDisplayed = action.visibility
    default:
        break
    }
    return state
}
