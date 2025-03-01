import { describe, expect, test } from "@jest/globals"
import type { Ingredient } from "./adjustSugar"
import { adjustSugarToIngredients } from "./adjustSugar"

type Scenario = {
  name: string
  ingredients: Ingredient[]
  percentage: number
  expectedSugar: number
}

describe("adjustSugarToIngredients", () => {
  const scenarios: Scenario[] = [
    {
      name: "Basic",
      ingredients: [
        {
          amount: 180,
          name: "flour",
          unit: "g",
        },
        {
          amount: 50,
          name: "sugar",
          unit: "g",
        },
      ],
      percentage: 10,
      expectedSugar: 20,
    },
    {
      name: "No sugar",
      ingredients: [
        {
          amount: 90,
          name: "flour",
          unit: "g",
        },
      ],
      percentage: 10,
      expectedSugar: 10,
    },
    {
      name: "Doesn't return multiple sugars",
      ingredients: [
        {
          amount: 90,
          name: "flour",
          unit: "g",
        },
        {
          amount: 90,
          name: "sugar",
          unit: "g",
        },
        {
          amount: 90,
          name: "cukor",
          unit: "g",
        },
      ],
      percentage: 10,
      expectedSugar: 10,
    },
    {
      name: "Converts the units",
      ingredients: [
        {
          amount: 20,
          name: "sugar",
          unit: "g",
        },
        {
          amount: 3,
          name: "flour",
          unit: "dkg",
        },
        {
          amount: 2,
          name: "apple",
          unit: "kg",
        },
      ],
      percentage: 12,
      expectedSugar: 277,
    },
    {
      name: "Converts the units back",
      ingredients: [
        {
          amount: 30,
          name: "sugar",
          unit: "dkg",
        },
        {
          amount: 2,
          name: "apple",
          unit: "kg",
        },
      ],
      percentage: 20,
      expectedSugar: 50,
    },
    {
      name: "Basic - German",
      ingredients: [
        {
          amount: 180,
          name: "flour",
          unit: "g",
        },
        {
          amount: 50,
          name: "Zucker",
          unit: "g",
        },
      ],
      percentage: 10,
      expectedSugar: 20,
    },
    {
      name: "Basic - Hungarian",
      ingredients: [
        {
          amount: 180,
          name: "flour",
          unit: "g",
        },
        {
          amount: 50,
          name: "CUKOR",
          unit: "g",
        },
      ],
      percentage: 10,
      expectedSugar: 20,
    },
  ]

  for (const scenario of scenarios) {
    test(scenario.name, () => {
      // when
      const result = adjustSugarToIngredients(
        scenario.ingredients,
        scenario.percentage,
      )
      expect(result.amount).toBe(scenario.expectedSugar)
    })
  }

  test("Original ingredients are not being modified", () => {
    // given
    const originalAmount = 100
    const ingredients: Ingredient[] = [
      {
        amount: 100,
        name: "sugar",
        unit: "g",
      },
    ]

    // when
    adjustSugarToIngredients(ingredients, 10)

    // then
    expect(ingredients[0].amount).toBe(originalAmount)
  })
})
