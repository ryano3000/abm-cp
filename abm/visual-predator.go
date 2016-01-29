package abm

import (
	"bytes"
	"errors"
	"fmt"
	"math"
	"math/rand"
	"sort"

	"github.com/benjamin-rood/abm-colour-polymorphism/calc"
	"github.com/benjamin-rood/abm-colour-polymorphism/colour"
	"github.com/benjamin-rood/abm-colour-polymorphism/geometry"
	"github.com/benjamin-rood/abm-colour-polymorphism/render"
)

// VisualPredator - Predator agent type for Predator-Prey ABM
type VisualPredator struct {
	pos           geometry.Vector //	position in the environment
	movS          float64         //	speed	/ movement range per turn
	movA          float64         //	acceleration
	tr            float64         // turn rate / range (in radians)
	dir           geometry.Vector //	must be implemented as a unit vector
	𝚯             float64         //	 heading angle
	lifespan      int
	hunger        int     //	counter for interval between needing food
	attackSuccess bool    //	if during the turn, the VP agent successfully ate a CP prey agent
	fertility     int     //	counter for interval between birth and sex
	gravid        bool    //	i.e. pregnant
	vsr           float64 //	visual search range
	γ             float64 //	visual acuity (as a standard, use 1.0) - this should be defined at the Environmental level?
	colImprint    colour.RGB
}

func vpTesterAgent(xPos float64, yPos float64) (tester VisualPredator) {
	tester = vpTestPop(1)[0]
	tester.pos[x] = xPos
	tester.pos[y] = yPos
	return
}

func vpTestPop(size int) []VisualPredator {
	return GeneratePopulationVP(size, TestContext)
}

// GeneratePopulationVP will create `size` number of Visual Predator agents
func GeneratePopulationVP(size int, context Context) (pop []VisualPredator) {
	for i := 0; i < size; i++ {
		agent := VisualPredator{}
		agent.pos = geometry.RandVector(context.Bounds)
		if context.VpAgeing {
			if context.RandomAges {
				agent.lifespan = calc.RandIntIn(int(float64(context.VpLifespan)*0.7), int(float64(context.VpLifespan)*1.3))
			} else {
				agent.lifespan = context.VpLifespan
			}
		} else {
			agent.lifespan = 99999
		}
		agent.movS = context.VpMovS
		agent.movA = context.VpMovA
		agent.𝚯 = rand.Float64() * (2 * math.Pi)
		agent.dir = geometry.UnitVector(agent.𝚯)
		agent.tr = context.VpTurn
		agent.vsr = context.Vsr
		agent.γ = context.Vγ //	visual acuity
		agent.hunger = 0
		agent.fertility = 1
		agent.gravid = false
		agent.colImprint = colour.RandRGB()
		pop = append(pop, agent)
	}
	return
}

// GetDrawInfo exports the data set needed for agent visualisation.
func (vp *VisualPredator) GetDrawInfo() (ar render.AgentRender) {
	ar.Type = "vp"
	ar.X = vp.pos[x]
	ar.Y = vp.pos[y]
	ar.Heading = vp.𝚯
	if vp.attackSuccess {
		ar.Colour = colour.RGB256{Red: 0, Green: 0, Blue: 0} // blink black on successful attack!
	} else {
		ar.Colour = vp.colImprint.To256()
	}
	return
}

// Turn updates 𝚯 and dir vector to the new heading offset by 𝚯
func (vp *VisualPredator) Turn(𝚯 float64) {
	newHeading := geometry.UnitAngle(vp.𝚯 + 𝚯)
	vp.dir[x] = math.Cos(newHeading)
	vp.dir[y] = math.Sin(newHeading)
	vp.𝚯 = newHeading
}

// Move updates the agent's position if it doesn't encounter any errors.
func (vp *VisualPredator) Move() error {
	var posOffset, newPos geometry.Vector
	var err error
	posOffset, err = geometry.VecScalarMultiply(vp.dir, vp.movS*vp.movA)
	if err != nil {
		return errors.New("agent move failed: " + err.Error())
	}
	newPos, err = geometry.VecAddition(vp.pos, posOffset)
	if err != nil {
		return errors.New("agent move failed: " + err.Error())
	}
	newPos[x] = calc.WrapFloatIn(newPos[x], -1.0, 1.0)
	newPos[y] = calc.WrapFloatIn(newPos[y], -1.0, 1.0)
	vp.pos = newPos
	return nil
}

// VSRSectorSampling checks which sectors the VP agent's
// Visual Search Radius intersects.
// This initial version samples from 4 points on the circumference
// of the circle with radius vp.visRange originating at the VP agent's position
// The four sample points on the circumference at 45°, 135°, 225°, 315°
// or π/4, 3π/4, 5π/4, 7π/4 radians,
// or NE, NW, SW, SE on a compass, if you want to think of it that way :-)
func (vp *VisualPredator) VSRSectorSampling(d float64, n int) ([4][2]int, error) {
	sectorSamples := [4][2]int{}

	x45 := vp.pos[x] + (vp.vsr * (math.Cos(math.Pi / 4)))
	y45 := vp.pos[y] + (vp.vsr * (math.Sin(math.Pi / 4)))

	x135 := vp.pos[x] + (vp.vsr * (math.Cos(3 * math.Pi / 4)))
	y135 := vp.pos[y] + (vp.vsr * (math.Sin(3 * math.Pi / 4)))

	x225 := vp.pos[x] + (vp.vsr * (math.Cos(5 * math.Pi / 4)))
	y225 := vp.pos[y] + (vp.vsr * (math.Sin(5 * math.Pi / 4)))

	x315 := vp.pos[x] + (vp.vsr * (math.Cos(7 * math.Pi / 4)))
	y315 := vp.pos[y] + (vp.vsr * (math.Sin(7 * math.Pi / 4)))

	sectorSamples[0][0], sectorSamples[0][1] = geometry.TranslatePositionToSector2D(d, n, geometry.Vector{x45, y45})

	sectorSamples[1][0], sectorSamples[1][1] = geometry.TranslatePositionToSector2D(d, n, geometry.Vector{x135, y135})

	sectorSamples[2][0], sectorSamples[2][1] = geometry.TranslatePositionToSector2D(d, n, geometry.Vector{x225, y225})

	sectorSamples[3][0], sectorSamples[3][1] = geometry.TranslatePositionToSector2D(d, n, geometry.Vector{x315, y315})

	return sectorSamples, nil
}

// PreySearch – uses Visual Search to try to 'recognise' a nearby prey agent within model Environment to target
func (vp *VisualPredator) PreySearch(prey []ColourPolymorphicPrey, searchChance float64) (target *ColourPolymorphicPrey, err error) {
	_ = "breakpoint" // godebug

	var searchSet []*ColourPolymorphicPrey
	for i := range prey { //	exhaustive search 😱
		prey[i].δ, err = geometry.VectorDistance(vp.pos, prey[i].pos)
		if prey[i].δ <= vp.vsr { // ∴ only include the prey agent for considertion if within visual range
			prey[i].𝛘 = colour.RGBDistance(vp.colImprint, prey[i].colouration)
			searchSet = append(searchSet, &prey[i])
		}
	}
	/*
		fmt.Println("Before sorting")
		for i := range prey {
			fmt.Printf("%v %v %v %v %v %p\n", i, prey[i].pos, prey[i].δ, prey[i].colouration, prey[i].𝛘, &prey[i])
		}
	*/
	sort.Sort(VisualDifferentiation(searchSet))
	/*
		// fmt.Println("After sorting")
		// for i := range searchSet {
		// 	fmt.Printf("%v %v %v %v %v %p\n", i, searchSet[i].pos, searchSet[i].δ, searchSet[i].colouration, searchSet[i].𝛘, searchSet[i])
		// }
	*/
	// search within biased and reduced set
	for i, p := range searchSet {
		if ((1.0 - p.𝛘) * (1.0 - p.δ) * vp.γ) > (1.0 - searchChance) {
			return searchSet[i], err
		}
	}
	return
}

// Intercept attempts to turn and move towards target position (as much as vp is able):
func (vp *VisualPredator) Intercept(prey *ColourPolymorphicPrey) (attack bool, err error) {
	if prey == nil {
		return
	}
	vx := prey.pos
	dist := prey.δ
	Ψ, err := geometry.AngleToIntercept(vp.pos, vp.𝚯, vx)
	if dist < vp.movS {
		attack = true
		vp.pos = vx
		vp.Turn(Ψ)
		return
	}
	vp.Turn(calc.ClampFloatIn(Ψ, -vp.tr, vp.tr))
	vp.Move()
	return
}

// Attack VP agent attempts to attack CP prey agent
func (vp *VisualPredator) Attack(prey *ColourPolymorphicPrey, vpAttackChance float64, imprintFactor float64) {
	if prey == nil {
		return
	}
	vpAttackChance = 1 - vpAttackChance
	α := rand.Float64()
	if α > vpAttackChance {
		vp.colourImprinting(prey.colouration, imprintFactor)
		vp.hunger -= 5
		prey.lifespan = 0 //	i.e. prey agent flagged for removal at the beginning of next turn and will not be drawn again.
		vp.attackSuccess = true
		fmt.Println("eaten =", prey.String())
		fmt.Println("eater =", vp.String())
	}
}

// colourImprinting updates VP colour / visual recognition bias
// Uses a bias / weighting value, 𝜎 (sigma) to control the degree of
// adaptation VP will make to differences in 'eaten' CPP colours.
func (vp *VisualPredator) colourImprinting(target colour.RGB, 𝜎 float64) error {
	𝚫red := (vp.colImprint.Red - target.Red) * 𝜎
	𝚫green := (vp.colImprint.Green - target.Green) * 𝜎
	𝚫blue := (vp.colImprint.Blue - target.Blue) * 𝜎
	vp.colImprint.Red = vp.colImprint.Red - 𝚫red
	vp.colImprint.Green = vp.colImprint.Green - 𝚫green
	vp.colImprint.Blue = vp.colImprint.Blue - 𝚫blue
	return nil
}

// animal-agent Mortal interface methods:

// Age the vp agent
func (vp *VisualPredator) Age(ctxt Context) string {
	vp.attackSuccess = false
	vp.hunger++
	if ctxt.VpAgeing {
		vp.lifespan--
	}
	return vp.jump(ctxt)
}

func (vp *VisualPredator) jump(ctxt Context) (jump string) {
	switch {
	case vp.lifespan <= 0:
		jump = "DEATH"
	case vp.hunger > ctxt.VpHungerLimit && ctxt.Starvation:
		jump = "DEATH"
	case vp.hunger < ctxt.VpSexualRequirement:
		jump = "MATE SEARCH"
	default:
		jump = "PREY SEARCH"
	}
	return
}

// String returns a clear textual presentation the internal values of the VP agent
func (vp *VisualPredator) String() string {
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("pos=(%v,%v)\n", vp.pos[x], vp.pos[y]))
	buffer.WriteString(fmt.Sprintf("movS=%v\n", vp.movS))
	buffer.WriteString(fmt.Sprintf("movA=%v\n", vp.movA))
	buffer.WriteString(fmt.Sprintf("𝚯=%v\n", vp.𝚯))
	buffer.WriteString(fmt.Sprintf("dir=(%v,%v)\n", vp.dir[x], vp.dir[y]))
	buffer.WriteString(fmt.Sprintf("tr=%v\n", vp.tr))
	buffer.WriteString(fmt.Sprintf("Vsr=%v\n", vp.vsr))
	buffer.WriteString(fmt.Sprintf("lifespan=%v\n", vp.lifespan))
	buffer.WriteString(fmt.Sprintf("hunger=%v\n", vp.hunger))
	buffer.WriteString(fmt.Sprintf("fertility=%v\n", vp.fertility))
	buffer.WriteString(fmt.Sprintf("gravid=%v\n", vp.gravid))
	buffer.WriteString(fmt.Sprintf("colouration=%v\n", vp.colImprint))
	return buffer.String()
}
