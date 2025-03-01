const conversionMap = {
  g: 1,
  dkg: 10,
  kg: 1000,
  oz: 28.3,
  lb: 453.59237,
}

export type Unit = keyof typeof conversionMap

const defaultUnit: Unit = "g"

const sugarNames = ["sugar", "zucker", "cukor"]

export type Ingredient = {
  unit: Unit
  amount: number
  name: string
}

export function adjustSugarToIngredients(
  ingredients: Ingredient[],
  sugarPercentage: number,
): Ingredient {
  const sum = sumUpNotSugar(ingredients)
  const ratio = sugarPercentage / 100
  const newSugarAmount = (sum.amount * ratio) / (1 - ratio)

  for (const ingredient of ingredients) {
    if (!isSugar(ingredient.name)) {
      continue
    }

    const sugar: Ingredient = { ...ingredient }
    return {
      name: ingredient.name,
      unit: ingredient.unit,
      amount: Math.round(newSugarAmount * (1 / conversionMap[sugar.unit])),
    }
  }

  return {
    name: "Sugar",
    amount: Math.round(newSugarAmount),
    unit: defaultUnit,
  }
}

function sumUpNotSugar(ingredients: Ingredient[]): Ingredient {
  const sum: Ingredient = {
    name: "everything excep sugar",
    unit: defaultUnit,
    amount: 0,
  }
  for (const ingredient of ingredients) {
    if (!isSugar(ingredient.name)) {
      sum.amount += ingredient.amount * conversionMap[ingredient.unit]
    }
  }

  return sum
}

export function isSugar(name: string): boolean {
  return sugarNames.includes(name.toLocaleLowerCase().trim())
}
