<template>
  <div class="flex-col text-center">
    <h1>Adjust the sugar in recipes</h1>
    <p>Do you find the recipes contain too much sugar?</p>
    <p>Here's a tool which helps adjusting it!</p>
    <p class="mt-5">Add the original ingredients:</p>
    <form
      @submit.prevent="addIngredient"
      class="grid grid-flow-row lg:grid-flow-col justify-center mt-4 mb-10"
    >
      <Input
        id="input-amount"
        label="Amount"
        type="number"
        v-model="state.ingredientAmount"
        required
        :step="0.01"
        class="mr-2 mt-3 col-span-2"
      />
      <Select
        id="select-unit"
        v-model="state.ingredientUnit"
        :options="['g', 'dkg', 'kg', 'oz', 'lb']"
        label="Unit"
        class="mr-2 mt-3 col-span-1"
      />
      <Input
        id="input-name"
        label="Name"
        type="text"
        v-model="state.ingredientName"
        class="mr-2 mt-3 col-span-4"
        required
      />
      <div class="col-span-4 flex">
        <button
          type="submit"
          class="lg:self-end mx-auto lg:mt-0 mt-3 bg-primary lg:my-0 h-10 px-2 text-white rounded-lg shadow-md hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-opacity-75"
        >
          Add ingredient
        </button>
      </div>
    </form>
    <IngredientComponent
      v-for="ingredient in state.ingredients"
      v-bind="ingredient"
      class="mt-1"
    />
    <div v-if="ingredientsExist" class="flex justify-center mt-4 mb-10">
      <Input
        id="input-sugar"
        label="Wanted sugar percentage in the whole recipe"
        type="number"
        v-model="state.sugarPercentage"
        required
        :step="1"
        :min="1"
        :max="99"
        class="mr-2"
      />
    </div>
    <p v-if="ingredientsExist" class="text-xl">
      The adjusted sugar amount:
      <strong>{{ `${adjustedSugar.amount} ${adjustedSugar.unit}` }}</strong>
    </p>
  </div>
</template>

<script setup lang="ts">
import { computed, reactive } from "vue"
import type { Props as IngredientProps } from "./Ingredient.vue"
import IngredientComponent from "./Ingredient.vue"
import Input from "./Input.vue"
import {
  adjustSugarToIngredients,
  type Ingredient,
  type Unit,
} from "./adjustSugar"
import Select from "./Select.vue"

type IngredientMap = {
  [name: string]: IngredientProps
}

type State = {
  ingredients: IngredientMap
  adjustedIngredients: IngredientProps[]
  ingredientName: string
  ingredientAmount?: number
  ingredientUnit: Unit
  sugarPercentage: number
}

const adjustedSugar = computed<IngredientProps>(() => {
  const logicalIngredients = Object.values(state.ingredients).map(
    ingredientPropToIngredient,
  )
  return adjustSugarToIngredients(logicalIngredients, state.sugarPercentage)
})

function ingredientPropToIngredient(prop: IngredientProps): Ingredient {
  return {
    name: prop.name,
    unit: prop.unit,
    amount: prop.amount || 0,
  }
}

const state = reactive<State>({
  ingredients: {},
  adjustedIngredients: [],
  ingredientName: "",
  ingredientAmount: undefined,
  ingredientUnit: "g",
  sugarPercentage: 16,
})

const ingredientsExist = computed<boolean>(() => {
  return Object.keys(state.ingredients).length > 0
})

function addIngredient() {
  const key = `${Date.now()}`

  state.ingredients[key] = {
    name: state.ingredientName,
    amount: state.ingredientAmount,
    unit: state.ingredientUnit,
    remove: () => {
      delete state.ingredients[key]
    },
  }

  state.ingredientName = ""
  state.ingredientAmount = undefined
}
</script>
