<script lang="ts">
	import loader from '@monaco-editor/loader';
	import type MonacoE from 'monaco-editor';
	import { onDestroy, onMount } from 'svelte';

	type Props = {
		editor?: MonacoE.editor.IStandaloneCodeEditor;
		value: string;
		options?: MonacoE.editor.IStandaloneEditorConstructionOptions;
		monaco?: typeof MonacoE;
	};

	let {
		editor = $bindable(),
		value = $bindable(),
		options = { value, automaticLayout: true },
		monaco = $bindable()
	}: Props = $props();

	let container: HTMLDivElement | undefined = $state();

	onMount(async () => {
		monaco = await loader.init();
		editor = monaco?.editor.create(container!, options);

		editor?.getModel()!.onDidChangeContent(() => {
			if (!editor) {
				return;
			}

			value = editor.getValue();
		});
	});

	onDestroy(() => editor?.dispose());
</script>

<div class="monaco-editor m-0 h-full w-full rounded-lg border p-0" bind:this={container}></div>
