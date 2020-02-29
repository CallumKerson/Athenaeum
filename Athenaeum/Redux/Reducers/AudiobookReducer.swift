/**
 AudiobookReducer.swift
 Copyright (c) 2020 Callum Kerr-Edwards
 */

import Foundation

func audiobookStateReducer(state: AudiobookState, action: Action) -> AudiobookState {
    var state = state
    switch action {
    case let action as AudiobookActions.SetAudiobook:
        state.audiobooks[action.audiobook.id] = action.audiobook
        state.importsInProgress.remove(action.audiobook.id)
    case let action as AudiobookActions.StartingImportOfAudiobook:
        state.importsInProgress.insert(action.id)
    case let action as ErrorActions.SetImportedFileAlreadyExistsError:
        state.importsInProgress.remove(action.idNotProcessed)
    case let action as ErrorActions.SetImportedFileIsOfWrongTypeError:
        state.importsInProgress.remove(action.idNotProcessed)
    case let action as ErrorActions.SetImportedFileURLCannotBeOpenedError:
        state.importsInProgress.remove(action.idNotProcessed)
    case let action as AudiobookActions.SetSelectedAudiobook:
        state.selectedAudiobook = action.audiobook
    default:
        break
    }
    return state
}
