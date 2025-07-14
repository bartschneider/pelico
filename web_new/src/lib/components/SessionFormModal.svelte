<script lang="ts">
  import { createEventDispatcher, onMount } from 'svelte';
  import type { Game, PlaySession } from '$lib/models';
  
  export let show = false;
  export let game: Game;
  export let session: PlaySession | null = null;
  
  const dispatch = createEventDispatcher();
  
  // Form fields
  let startTime = '';
  let endTime = '';
  let rating = '';
  let notes = '';
  
  // Initialize default values on mount
  onMount(() => {
    const now = new Date();
    startTime = now.toISOString().slice(0, 16);
  });
  
  // Initialize/reset form when modal opens
  function initializeForm() {
    if (session) {
      // Edit mode
      startTime = new Date(session.start_time).toISOString().slice(0, 16);
      endTime = session.end_time ? new Date(session.end_time).toISOString().slice(0, 16) : '';
      rating = session.rating?.toString() || '';
      notes = session.notes || '';
    } else {
      // Create mode - default to current time
      const now = new Date();
      startTime = now.toISOString().slice(0, 16);
      endTime = '';
      rating = '';
      notes = '';
    }
  }
  
  // Call initializeForm when show changes to true
  $: if (show) {
    initializeForm();
  }
  
  function handleSubmit() {
    if (!startTime) {
      alert('Please select a start time');
      return;
    }
    
    try {
      const sessionData: Partial<PlaySession> = {
        start_time: new Date(startTime).toISOString(),
        end_time: endTime ? new Date(endTime).toISOString() : undefined,
        rating: rating ? parseInt(rating) : undefined,
        notes: notes || undefined
      };
      
      dispatch('submit', sessionData);
    } catch (error) {
      console.error('Error creating session data:', error);
      alert('Error preparing session data. Please check your inputs.');
    }
  }
  
  function handleClose() {
    dispatch('close');
  }
  
  function handleBackdropClick(event: MouseEvent) {
    if (event.target === event.currentTarget) {
      handleClose();
    }
  }
</script>

{#if show}
  <div class="modal fade show d-block" tabindex="-1" role="dialog" on:click={handleBackdropClick}>
    <div class="modal-dialog modal-dialog-centered" role="document">
      <div class="modal-content">
        <div class="modal-header">
          <h5 class="modal-title">
            {session ? 'Edit' : 'Log'} Play Session - {game.title}
          </h5>
          <button type="button" class="btn-close" aria-label="Close" on:click={handleClose}></button>
        </div>
        <form on:submit|preventDefault={handleSubmit}>
          <div class="modal-body">
            <div class="mb-3">
              <label for="startTime" class="form-label">Start Time <span class="text-danger">*</span></label>
              <input 
                type="datetime-local" 
                class="form-control" 
                id="startTime" 
                bind:value={startTime}
                required
              />
            </div>
            
            <div class="mb-3">
              <label for="endTime" class="form-label">End Time</label>
              <input 
                type="datetime-local" 
                class="form-control" 
                id="endTime" 
                bind:value={endTime}
                min={startTime}
              />
              {#if !endTime}
                <small class="text-muted">Leave empty for ongoing session</small>
              {/if}
            </div>
            
            <div class="mb-3">
              <label for="rating" class="form-label">Rating (1-10)</label>
              <input 
                type="number" 
                class="form-control" 
                id="rating" 
                bind:value={rating}
                min="1"
                max="10"
                placeholder="How was this session?"
              />
            </div>
            
            <div class="mb-3">
              <label for="notes" class="form-label">Notes</label>
              <textarea 
                class="form-control" 
                id="notes" 
                rows="3"
                bind:value={notes}
                placeholder="Any thoughts about this play session?"
                maxlength="1000"
              ></textarea>
              <small class="text-muted">{notes.length}/1000 characters</small>
            </div>
            
            {#if startTime && endTime}
              <div class="alert alert-info">
                <i class="fas fa-clock me-2"></i>
                Duration: {calculateDuration(startTime, endTime)}
              </div>
            {/if}
          </div>
          <div class="modal-footer">
            <button type="button" class="btn btn-secondary" on:click={handleClose}>Cancel</button>
            <button type="submit" class="btn btn-primary">
              <i class="fas fa-save me-1"></i>
              {session ? 'Update' : 'Log'} Session
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
  <div class="modal-backdrop fade show"></div>
{/if}

<style>
  .modal {
    display: block;
  }
  
  .modal-backdrop {
    position: fixed;
    top: 0;
    left: 0;
    z-index: 1040;
    width: 100vw;
    height: 100vh;
    background-color: #000;
    opacity: 0.5;
  }
  
  .modal-dialog {
    z-index: 1050;
  }
</style>

<script lang="ts" context="module">
  function calculateDuration(start: string, end: string): string {
    const startDate = new Date(start);
    const endDate = new Date(end);
    const diffMs = endDate.getTime() - startDate.getTime();
    const diffMins = Math.floor(diffMs / 60000);
    
    if (diffMins < 60) {
      return `${diffMins} minute${diffMins !== 1 ? 's' : ''}`;
    }
    
    const hours = Math.floor(diffMins / 60);
    const mins = diffMins % 60;
    
    if (mins === 0) {
      return `${hours} hour${hours !== 1 ? 's' : ''}`;
    }
    
    return `${hours} hour${hours !== 1 ? 's' : ''} ${mins} minute${mins !== 1 ? 's' : ''}`;
  }
</script>