# Nano Editor Feature Plan

## Overview
Implement a new external prompt editor feature that allows users to switch between basic and nano editor for composing prompts. When activated with ctrl+n, nano opens in a full popup window that is smaller than the CLI window and uses the same popup motor as other settings dialogs.

## Current State Analysis
- **Existing External Editor**: Already supports ctrl+e to open $EDITOR (defaults to nvim/notepad)
- **Dialog System**: Well-established with DialogModel interface and dialogs package
- **Editor Component**: Located in `internal/tui/components/chat/editor/`
- **Key Bindings**: Currently ctrl+o for external editor

## Feature Requirements
1. **Editor Selection**: Allow switching between basic textarea and nano editor
2. **Nano Popup**: Full popup window smaller than CLI window
3. **Same Popup Motor**: Use existing dialogs system like ctrl+p menu
4. **Nano Integration**: Display actual nano software in popup
5. **Seamless Workflow**: User edits in nano, closes, popup disappears
6. **Auto-execution**: App takes prompt without modifications and runs immediately

## Implementation Plan

### Phase 1: Core Infrastructure

#### 1.1 Create Nano Dialog Component
**Location**: `internal/tui/components/dialogs/nano/`
**Files to Create**:
- `nano.go` - Main nano dialog implementation
- `keys.go` - Key bindings for nano dialog
- `nano_dialog.go` - Dialog model interface implementation

**Responsibilities**:
- Implement DialogModel interface
- Handle nano process lifecycle
- Manage popup positioning and sizing
- Handle content transfer back to editor

#### 1.2 Extend Editor KeyMap
**File**: `internal/tui/components/chat/editor/keys.go`
**Changes**:
- Add `OpenNanoEditor` key binding (ctrl+n)
- Update KeyBindings() method to include new binding

#### 1.3 Update Editor Component
**File**: `internal/tui/components/chat/editor/editor.go`
**Changes**:
- Add `openNanoEditor()` method
- Handle `OpenNanoEditorMsg` message
- Integrate with nano dialog system

### Phase 2: Nano Process Management

#### 2.1 Nano Process Handler
**Implementation Details**:
- Create temporary file for nano to edit
- Launch nano with proper terminal integration
- Monitor nano process completion
- Read back edited content
- Clean up temporary files

#### 2.2 Terminal Integration
**Key Challenges**:
- Ensure nano runs properly within TUI environment
- Handle stdin/stdout/stderr redirection
- Maintain terminal state during nano execution
- Restore TUI state after nano exits

#### 2.3 Content Transfer
**Mechanism**:
- Write current prompt content to temp file
- Launch nano with temp file
- Read content back after nano exits
- Return content via message system
- Auto-execute prompt immediately

### Phase 3: Dialog System Integration

#### 3.1 Dialog Registration
**File**: `internal/tui/components/dialogs/dialogs.go`
**Changes**:
- Add nano dialog ID constant
- Register nano dialog in dialog manager

#### 3.2 Message Flow
**New Messages**:
- `OpenNanoEditorMsg` - Trigger nano editor dialog
- `NanoEditorResultMsg` - Return edited content
- `NanoEditorClosedMsg` - Handle dialog closure

#### 3.3 Dialog Lifecycle
**States**:
- Dialog opens with current prompt content
- User edits in nano
- Dialog closes on nano exit
- Content automatically sent to chat

### Phase 4: UI/UX Enhancements

#### 4.1 Popup Styling
**Requirements**:
- Smaller than CLI window (80% width/height)
- Centered positioning
- Consistent styling with other dialogs
- Proper border and background

#### 4.2 Loading States
**User Feedback**:
- Show loading indicator while nano starts
- Display status during nano editing
- Handle nano startup errors gracefully

#### 4.3 Error Handling
**Scenarios**:
- Nano not installed
- Temporary file creation failure
- Nano process errors
- Content read/write failures

### Phase 5: Configuration and Preferences

#### 5.1 Editor Mode Settings
**Options**:
- Default editor mode (basic/nano)
- Nano-specific preferences
- Remember user choice

#### 5.2 Fallback Behavior
**Logic**:
- If nano not available, fallback to existing external editor
- Graceful degradation with user notification
- Configurable alternative editors

### Phase 6: Testing and Validation

#### 6.1 Unit Tests
**Coverage**:
- Nano dialog creation and lifecycle
- Process management functions
- Content transfer mechanisms
- Error handling scenarios

#### 6.2 Integration Tests
**Scenarios**:
- Full workflow from ctrl+n to prompt execution
- Dialog positioning and sizing
- Content preservation during editing
- Error recovery and user feedback

#### 6.3 Cross-Platform Testing
**Platforms**:
- macOS (primary target)
- Linux
- Windows (if nano available)

## Technical Challenges

### 1. Terminal Multiplexing
**Issue**: Running nano within a TUI application requires careful terminal state management
**Solution**: Use proper process execution with stdin/stdout redirection and terminal restoration

### 2. Process Synchronization
**Issue**: Coordinating between TUI event loop and nano process lifecycle
**Solution**: Use tea.ExecProcess with proper callback handling

### 3. Content Integrity
**Issue**: Ensuring edited content is properly transferred back to the application
**Solution**: Atomic file operations with validation and error handling

### 4. User Experience Consistency
**Issue**: Maintaining TUI responsiveness during external editor usage
**Solution**: Proper loading states and user feedback throughout the process

## Dependencies

### External Dependencies
- `nano` editor must be installed on system
- Proper terminal environment for nano execution

### Internal Dependencies
- Existing dialogs system
- Editor component infrastructure
- Message passing system
- File management utilities

## Success Criteria

1. **Functionality**: Ctrl+n opens nano in popup, editing works seamlessly
2. **Performance**: No noticeable lag in popup opening/closing
3. **Reliability**: Proper error handling for edge cases
4. **UX**: Intuitive workflow with clear feedback
5. **Integration**: Works with existing TUI components and dialogs

## Implementation Timeline

### Week 1: Core Infrastructure
- Create nano dialog component
- Implement basic key bindings
- Set up process management foundation

### Week 2: Process Integration
- Implement nano execution logic
- Handle terminal state management
- Basic content transfer mechanism

### Week 3: Dialog Integration
- Integrate with existing dialog system
- Implement proper message flow
- Add loading states and error handling

### Week 4: Polish and Testing
- UI/UX refinements
- Comprehensive testing
- Documentation and cleanup

## Risk Assessment

### High Risk
- Terminal multiplexing issues
- Cross-platform compatibility
- Process lifecycle management

### Medium Risk
- Content transfer reliability
- User experience consistency
- Error handling completeness

### Low Risk
- Dialog system integration
- Key binding conflicts
- Configuration management

## Alternative Approaches

### 1. Embedded Nano
**Pros**: Better integration, no external process
**Cons**: Complex implementation, potential licensing issues

### 2. Web-based Editor
**Pros**: Cross-platform, rich features
**Cons**: Requires web server, changes user workflow

### 3. Enhanced Textarea
**Pros**: Simple, no external dependencies
**Cons**: Limited editing capabilities compared to nano

## Conclusion

This feature will significantly enhance the prompt editing experience by providing users with a familiar, powerful text editor while maintaining the benefits of the TUI interface. The implementation leverages existing infrastructure and follows established patterns in the codebase.</content>
<parameter name="file_path">develop/nano_editor_feature.md