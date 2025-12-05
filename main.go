package main

import (
    "encoding/json"
    "fmt"
    "log"
    "math"
    "net/http"
    "os"
    "time"
    swe "github.com/astrotools/swephgo/swe"

)

// Structs for input and output
type BirthData struct {
    Year     int     `json:"year"`
    Month    int     `json:"month"`
    Day      int     `json:"day"`
    Hour     float64 `json:"hour"`
    Lat      float64 `json:"lat"`
    Lon      float64 `json:"lon"`
    Timezone string  `json:"timezone"`
}

type ChironReading struct {
    Sign             string  `json:"sign"`
    Degree           float64 `json:"degree"`
    House            int     `json:"house"`
    TraditionalWound string  `json:"traditional_wound"`
    LHPStrength      string  `json:"lhp_strength"`
    Timestamp        int64   `json:"timestamp"`
}

func main() {
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    http.HandleFunc("/", homeHandler)
    http.HandleFunc("/api/health", healthHandler)
    http.HandleFunc("/api/chiron", chironHandler)

    log.Printf("🚀 Chiron Oracle starting on port %s", port)
    log.Printf("📡 Local: http://localhost:%s", port)
    if err := http.ListenAndServe(":"+port, nil); err != nil {
        log.Fatal(err)
    }
    
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path != "/" {
        http.NotFound(w, r)
        return
        // Julian Day from a UTC time
func julianDay(t time.Time) float64 {
    year := t.Year()
    month := int(t.Month())
    day := t.Day()
    hour := float64(t.Hour()) + float64(t.Minute())/60.0 + float64(t.Second())/3600.0

    if month <= 2 {
        year -= 1
        month += 12
    }
    A := year / 100
    B := 2 - A + A/4

    jd := math.Floor(365.25*float64(year+4716)) +
        math.Floor(30.6001*float64(month+1)) +
        float64(day) + float64(B) - 1524.5 +
        hour/24.0

    return jd
}

// Tropical ecliptic longitude of Chiron via Swiss Ephemeris
func computeChironLongitude(jd float64) float64 {
    xx := make([]float64, 6)
    serr := make([]byte, 256)

    if ret := swe.CalcUt(jd, swe.SE_CHIRON, swe.SEFLG_SWIEPH, xx, serr); ret < 0 {
        log.Printf("Swiss Ephemeris error: %s", string(serr))
        return 0.0
    }
    return xx[0] // degrees, 0–360 (tropical)
}

// Map longitude to zodiac sign
func signFromLongitude(longDeg float64) string {
    signs := []string{"Aries", "Taurus", "Gemini", "Cancer", "Leo", "Virgo",
        "Libra", "Scorpio", "Sagittarius", "Capricorn", "Aquarius", "Pisces"}
    idx := int(math.Floor(longDeg / 30.0)) % 12
    return signs[idx]
}

// Whole-sign house mapping from ASC sign index
func wholeSignHouse(ascSignIndex int, longDeg float64) int {
    signIndex := int(math.Floor(longDeg / 30.0)) % 12
    if signIndex < 0 {
        signIndex += 12
    }
    dist := signIndex - ascSignIndex
    if dist < 0 {
        dist += 12
    }
    return dist + 1
}

    }
    html := `
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Chiron Wound Inversion Oracle</title>
  <style>
    body {
      font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', system-ui, sans-serif;
      background: linear-gradient(135deg, #020617, #0f172a);
      color: #f8fafc;
      min-height: 100vh;
      line-height: 1.6;
    }
    .container { max-width: 800px; margin: 0 auto; padding: 2rem; }
    h1 { font-size: 3rem; background: linear-gradient(45deg, #8b5cf6, #3b82f6);
         -webkit-background-clip: text; color: transparent; text-align: center; }
    .card { background: rgba(255,255,255,0.05); border-radius: 15px; padding: 2rem; margin-bottom: 2rem; }
    label { display:block; margin-bottom:0.5rem; color:#94a3b8; font-weight:500; }
    input { width:100%; padding:0.75rem; background:rgba(255,255,255,0.1); border:1px solid rgba(255,255,255,0.2);
            border-radius:10px; color:white; margin-bottom:1rem; }
    button { background: linear-gradient(45deg, #8b5cf6, #3b82f6); color:white; border:none; padding:1rem 2rem;
             border-radius:10px; font-size:1.1rem; font-weight:600; cursor:pointer; width:100%; }
    .result { background: rgba(255,255,255,0.05); border-radius: 15px; padding: 2rem; margin-top: 2rem; }
  </style>
</head>
<body>
  <div class="container">
    <h1>🔮 Chiron Wound Inversion Oracle</h1>
    <div class="card">
      <h2>📝 Birth Details</h2>
      <label for="year">Year</label>
      <input type="number" id="year" value="1990">
      <label for="month">Month</label>
      <input type="number" id="month" value="5">
      <label for="day">Day</label>
      <input type="number" id="day" value="12">
      <label for="hour">Hour (24h)</label>
      <input type="number" id="hour" value="14">
    </div>
    <div class="card">
      <h2>📍 Birth Location</h2>
      <label for="place">City, Country</label>
      <input type="text" id="place" placeholder="Kochi, India">
      <small style="color:#94a3b8;">Type your city and country, e.g. "London, UK"</small>
    </div>
    <button id="calculateBtn" onclick="getReading()">✨ Calculate My Chiron Reading</button>
    <div class="result" id="result"></div>
  </div>

  <script>
    async function getCoordinates(place) {
      const url = "https://nominatim.openstreetmap.org/search?format=json&q=" + encodeURIComponent(place);
      const response = await fetch(url);
      const data = await response.json();
      if (data.length > 0) {
        return { lat: parseFloat(data[0].lat), lon: parseFloat(data[0].lon) };
      } else {
        throw new Error("Place not found");
      }
    }

    async function getReading() {
      const btn = document.getElementById('calculateBtn');
      const resultDiv = document.getElementById('result');
      resultDiv.innerHTML = "⏳ Consulting the Oracle...";

      try {
        const place = document.getElementById('place').value;
        const coords = await getCoordinates(place);

        const data = {
          year: parseInt(document.getElementById('year').value),
          month: parseInt(document.getElementById('month').value),
          day: parseInt(document.getElementById('day').value),
          hour: parseFloat(document.getElementById('hour').value),
          lat: coords.lat,
          lon: coords.lon,
          timezone: "Asia/Kolkata" // TODO: auto-detect later
        };

        btn.disabled = true;
        const response = await fetch('/api/chiron', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify(data)
        });

        if (!response.ok) throw new Error("API Error: " + response.status);
        const reading = await response.json();

        resultDiv.innerHTML =
          '<h2>✨ Your Chiron Reading</h2>' +
          '<p><strong>Sign:</strong> ' + reading.sign + '</p>' +
          '<p><strong>Degree:</strong> ' + reading.degree + '°</p>' +
          '<p><strong>House:</strong> ' + reading.house + '</p>' +
          '<p><strong>Traditional Wound:</strong> ' + reading.traditional_wound + '</p>' +
          '<p><strong>LHP Strength:</strong> ' + reading.lhp_strength + '</p>';
      } catch (err) {
        resultDiv.innerHTML = "❌ Error: " + err.message;
      } finally {
        btn.disabled = false;
      }
    }
  </script>
</body>
</html>`
    w.Header().Set("Content-Type", "text/html")
    fmt.Fprint(w, html)
}
// --- Part 2: Backend Handlers + Calculation ---

func healthHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{
        "status":  "healthy",
        "service": "chiron-oracle",
        "version": "1.0.0",
        "time":    time.Now().Unix(),
    })
}
func chironHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

    if r.Method == "OPTIONS" {
        w.WriteHeader(http.StatusOK)
        return
    }

       var req BirthData
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    hour := int(req.Hour)
    minute := int((req.Hour - float64(hour)) * 60)

    loc, err := time.LoadLocation(req.Timezone)
    if err != nil {
        http.Error(w, "invalid timezone", http.StatusBadRequest)
        return
    }
    local := time.Date(req.Year, time.Month(req.Month), req.Day, hour, minute, 0, 0, loc)
    utc := local.UTC()

    jd := julianDay(utc)
    chironLon := computeChironLongitude(jd)

    sign := signFromLongitude(chironLon)
    degree := math.Mod(chironLon, 30)
    ascSignIndex := 7 // TODO: replace with real ASC calc
    house := wholeSignHouse(ascSignIndex, chironLon)

    resp := ChironReading{
        Sign:             sign,
        Degree:           math.Round(degree*100) / 100,
        House:            house,
        TraditionalWound: "TODO: fill in", // you can add logic later
        LHPStrength:      "TODO: fill in",
        Timestamp:        utc.Unix(),
    }

    json.NewEncoder(w).Encode(resp)
}



func calculateChiron(input BirthData) (string, float64, int) {
    seed := float64(input.Year*10000+input.Month*100+input.Day) +
        input.Hour + math.Abs(input.Lat) + math.Abs(input.Lon)

    signs := []string{
        "Aries", "Taurus", "Gemini", "Cancer", "Leo", "Virgo",
        "Libra", "Scorpio", "Sagittarius", "Capricorn", "Aquarius", "Pisces",
    }

    goldenRatio := 1.61803398875
    signIdx := int(math.Mod(math.Abs(seed)*goldenRatio, 12))
    degree := math.Mod(math.Abs(seed)*math.Pi, 30)

    houseSeed := math.Abs(input.Lat*input.Lon + input.Hour*100)
    house := int(math.Mod(houseSeed, 12)) + 1

    return signs[signIdx], math.Round(degree*100)/100, house
}
// --- Interpretations ---
func getInterpretation(sign string, house int) (string, string) {
    interpretations := map[string]map[int][2]string{
        "Aries": {
            1: {
                "You often struggle with impatience and uncertainty about who you truly are. This wound can manifest as frustration when trying to assert your identity or feeling like you must constantly prove yourself.",
                "Through inversion, this wound becomes a strength: you forge identity through fearless, decisive action. You learn to embrace the warrior within, carving out space for yourself without apology.",
            },
            2: {
                "Financial impulsiveness and insecurity may plague you, leading to rash decisions or fear of scarcity. You may feel wounded when your material stability is challenged.",
                "Inverted, this becomes bold initiative in wealth creation. You harness risk-taking as a tool, building resources through daring ventures and pioneering financial strategies.",
            },
            3: {
                "Communication can feel aggressive or misunderstood, leaving you wounded by conflict in dialogue. You may fear that your words are too sharp or unwelcome.",
                "Inverted, your fiery speech becomes persuasive conviction. You inspire others with passionate clarity, turning raw honesty into a force that moves minds and hearts.",
            },
            4: {
                "Family dynamics may wound you through struggles with leadership or dominance in the home. You may feel burdened by expectations or clashes with authority figures.",
                "Inverted, you channel this into leading your lineage into new territories. You become the pioneer who transforms ancestral wounds into fresh paths for growth.",
            },
            5: {
                "Creative impatience can leave you frustrated, feeling blocked or unable to sustain projects. You may fear that your creations lack depth or longevity.",
                "Inverted, you create from raw impulse and passion. Your art becomes a living flame, igniting inspiration in others through its immediacy and intensity.",
            },
            6: {
                "Workplace frustration and conflict may wound you, as you struggle with rigid systems or authority. You may feel trapped in routines that stifle your fire.",
                "Inverted, you break systems to rebuild them stronger. You become a catalyst for innovation, transforming frustration into revolutionary improvements.",
            },
            7: {
                "Partnerships may wound you through dominance struggles or fear of losing independence. You may feel torn between self-assertion and compromise.",
                "Inverted, you master the dance of equal opposition. You learn to wield strength in balance, forging partnerships that thrive on mutual respect and fiery passion.",
            },
            8: {
                "You may fear transformation, death, or loss of control, leaving you wounded by the intensity of change. This can manifest as resistance to deep psychological shifts.",
                "Inverted, you embrace the fires of transformation. You die and are reborn repeatedly, mastering the art of self-renewal and wielding power through fearless surrender.",
            },
            9: {
                "Philosophical certainty may wound you, leaving you rigid in beliefs or fearful of questioning dogma. You may cling to truths that feel safe but limiting.",
                "Inverted, you burn beliefs to find truth in ashes. You become a seeker who thrives in uncertainty, discovering wisdom through destruction and renewal of ideas.",
            },
            10: {
                "Career impatience may wound you, as you feel blocked by slow progress or external limitations. You may fear that success is always out of reach.",
                "Inverted, you achieve by ignoring conventional timelines. You blaze trails in public life, proving that ambition fueled by fire can defy all expectations.",
            },
           11: {
                "Social pioneering anxiety may wound you, leaving you fearful of rejection or isolation when trying to lead groups. You may feel misunderstood in collective settings.",
                "Inverted, you form tribes through shared fire. You become the spark that ignites communities, building circles that thrive on your bold vision.",
            },
            12: {
                "Spiritual impulsiveness may wound you, as you rush into mystical experiences without grounding. You may fear losing yourself in the unconscious.",
                "Inverted, you navigate the unconscious with warrior focus. You wield spiritual fire as a disciplined tool, mastering hidden realms with courage and clarity.",
            },
        }, // <-- closes Aries only

        "Taurus": {
            1: {
                "You may struggle with self-worth tied to material possessions or physical security. This wound can manifest as constant comparison, or feeling inadequate unless you have tangible proof of your value.",
                "Inverted, you realize your value is inherent and not dependent on external validation. You cultivate unshakable confidence, embodying stability and groundedness that inspires others.",
            },
            2: {
                "Possessiveness over resources may wound you, leading to fear of loss or clinging to what you own. This can create anxiety around money and material stability.",
                "Inverted, you accumulate resources as sacred ritual, not hoarding. You learn to steward wealth wisely, building abundance that nourishes both yourself and your community.",
            },
            3: {
                "Stubborn communication may wound you, making dialogue rigid or resistant to new ideas. You may feel unheard or unable to adapt in conversations.",
                "Inverted, you speak with the weight of earth itself. Your words carry authority and grounding, offering stability and clarity in chaotic times.",
            },
            4: {
                "Family security obsession may wound you, leaving you anxious about protecting loved ones or clinging to tradition. You may fear instability in your roots.",
                "Inverted, you build foundations that last centuries. You become the architect of enduring legacies, creating homes and families that thrive on resilience.",
            },
            5: {
                "Creative stagnation may wound you, leaving you fearful of change or hesitant to experiment. You may feel blocked when trying to express beauty.",
                "Inverted, you create beauty that anchors the soul. Your art becomes timeless, grounding others in serenity and harmony.",
            },
            6: {
                "Workplace rigidity may wound you, making you resistant to change or innovation. You may feel trapped in routines that stifle growth.",
                "Inverted, you build systems that endure all storms. Your persistence and reliability transform work into sacred service, ensuring stability for all.",
            },
            7: {
                "Relationship security fears may wound you, leaving you anxious about abandonment or overly dependent on stability. You may cling to partners for reassurance.",
                "Inverted, you treat partnership as sacred contract. You embody loyalty and devotion, creating bonds that are unbreakable and deeply nourishing.",
            },
            8: {
                "Transformation resistance may wound you, leaving you fearful of change or reluctant to let go. You may resist growth even when it is necessary.",
                "Inverted, you master the pace of change. You embrace transformation as natural, learning to evolve steadily and gracefully.",
            },
            9: {
                "Belief system materialism may wound you, tying faith to tangible proof or rejecting spirituality that feels abstract. You may struggle to trust the unseen.",
                "Inverted, you find divinity in tangible reality. You discover sacredness in the physical world, embodying spirituality through grounded practices.",
            },
            10: {
                "Career stability obsession may wound you, leaving you fearful of risk or overly attached to predictable paths. You may resist ambition that feels uncertain.",
                "Inverted, you build legacy through persistent effort. Your career becomes a monument to endurance, proving that slow and steady truly wins.",
            },
            11: {
                "Social value anxiety may wound you, leaving you fearful of rejection or questioning your worth in groups. You may feel invisible or undervalued.",
                "Inverted, your worth defines the circle. You become the anchor of communities, offering stability and reliability that others depend on.",
            },
            12: {
                "Spiritual materialism may wound you, leaving you clinging to rituals or physical symbols without deeper connection. You may fear the intangible.",
                "Inverted, you discover the divine in every atom. You embody spirituality through presence, grounding mystical truths in the physical world.",
            },
        },
        "Gemini": {
            1: {
                "You may struggle with identity fragmentation, feeling pulled in many directions or uncertain about who you truly are. This wound can manifest as restlessness or fear of being inconsistent.",
                "Inverted, you master multiplicity as a superpower. You embrace your many facets, showing others that identity can be fluid, adaptable, and endlessly creative.",
            },
            2: {
                "Scattered financial focus may wound you, leaving you anxious about money or unable to sustain stability. You may feel overwhelmed by too many options or inconsistent priorities.",
                "Inverted, you harness wealth through information flow. You become resourceful by connecting ideas, people, and opportunities, turning variety into abundance.",
            },
            3: {
                "Mental overload may wound you, leaving you exhausted by constant thoughts, ideas, and communication. You may fear that your mind is too chaotic to be useful.",
                "Inverted, you become the central switchboard. You thrive as a communicator, weaving connections and insights that others depend on for clarity and innovation.",
            },
            4: {
                "Family communication issues may wound you, leaving you feeling unheard or misunderstood in your lineage. You may struggle with expressing truth in close relationships.",
                "Inverted, you rewrite family narratives. You become the storyteller who heals ancestral wounds, bringing new language and perspective to old patterns.",
            },
            5: {
                "Creative dilettantism may wound you, leaving you fearful of being unfocused or superficial in your art. You may feel blocked by too many interests.",
                "Inverted, you synthesize all forms into new creation. Your creativity thrives on diversity, producing works that are eclectic, innovative, and alive.",
            },
            6: {
                "Workplace distraction may wound you, leaving you scattered or unable to finish tasks. You may fear being seen as unreliable.",
                "Inverted, you master multitasking as sacred art. You show that adaptability and variety are strengths, bringing flexibility and innovation to any environment.",
            },
            7: {
                "Relationship indecision may wound you, leaving you fearful of commitment or overwhelmed by choices. You may struggle to balance curiosity with stability.",
                "Inverted, you embrace every connection as a teacher. You learn from each relationship, weaving wisdom from diversity and keeping love alive through curiosity.",
            },
            8: {
                "Intellectualizing transformation may wound you, leaving you fearful of surrender or overly analytical about deep change. You may resist emotional depth.",
                "Inverted, you understand death to master rebirth. You bring clarity and language to transformation, guiding others through change with insight and perspective.",
            },
            9: {
                "Belief system confusion may wound you, leaving you fearful of committing to one truth or overwhelmed by contradictions. You may feel lost in endless questioning.",
                "Inverted, you embrace truth as multifaceted. You show that wisdom is not singular, but woven from many perspectives, making you a bridge between worlds.",
            },
            10: {
                "Career versatility anxiety may wound you, leaving you fearful of being unfocused or undervalued. You may struggle to choose one path.",
                "Inverted, you shape-shift through professional realms. You thrive in variety, proving that adaptability and curiosity are assets in any career.",
            },
            11: {
                "Social butterfly syndrome may wound you, leaving you fearful of being seen as shallow or inconsistent in groups. You may feel anxious about belonging.",
                "Inverted, you network as a neural web. You become the connector who brings people together, weaving communities through your endless curiosity and communication.",
            },
            12: {
                "Spiritual restlessness may wound you, leaving you fearful of stillness or unable to commit to one practice. You may feel scattered in your search for meaning.",
                "Inverted, you discover that the void speaks all languages. You thrive spiritually by embracing diversity, finding sacredness in multiplicity and constant exploration.",
            },
        },
        "Cancer": {
            1: {
                "You may struggle with emotional boundaries, feeling overwhelmed by the needs of others or uncertain where you end and they begin. This wound can manifest as vulnerability or fear of being consumed by relationships.",
                "Inverted, you learn to wield empathy as a strength. You nurture without losing yourself, becoming a source of emotional resilience and guidance for those around you.",
            },
            2: {
                "Emotional attachment to possessions may wound you, leaving you fearful of losing what anchors your heart. You may cling to material things as symbols of safety.",
                "Inverted, you transform resources into emotional anchors. You cultivate security through meaningful possessions, imbuing them with love and memory rather than fear.",
            },
            3: {
                "Communication sensitivity may wound you, leaving you fearful of criticism or easily hurt by words. You may struggle to express yourself without fear of rejection.",
                "Inverted, your words become healing instruments. You speak with compassion and emotional depth, offering language that soothes wounds and builds bridges.",
            },
            4: {
                "Family emotional baggage may wound you, leaving you burdened by ancestral pain or unresolved dynamics. You may feel trapped by lineage expectations.",
                "Inverted, you transform lineage pain into power. You become the healer of your family, rewriting narratives and creating emotional sanctuaries for future generations.",
            },
            5: {
                "Creative vulnerability may wound you, leaving you fearful of exposing your inner world. You may hesitate to share art that feels too raw.",
                "Inverted, you create from emotional truth. Your vulnerability becomes your greatest strength, inspiring others through authenticity and courage.",
            },
            6: {
                "Workplace emotional labor may wound you, leaving you drained by caretaking roles or undervalued for your sensitivity. You may feel exploited for your compassion.",
                "Inverted, you wield care as strategic advantage. You transform emotional labor into leadership, showing that empathy is a powerful force in professional spaces.",
            },
            7: {
                "Relationship dependency may wound you, leaving you fearful of abandonment or overly reliant on others for emotional stability. You may struggle with independence.",
                "Inverted, you embrace interdependence as strength. You build partnerships rooted in mutual care, proving that vulnerability can coexist with resilience.",
            },
            8: {
                "Psychological depth fears may wound you, leaving you hesitant to explore the unconscious or fearful of emotional intensity. You may resist transformation.",
                "Inverted, you navigate emotional underworlds with courage. You master hidden realms, turning fear into wisdom and emotional power.",
            },
            9: {
                "Belief system emotionality may wound you, leaving you fearful of faith that feels too vulnerable or dependent on feelings. You may struggle to trust intuition.",
                "Inverted, you root faith in feeling. You embrace emotional wisdom as sacred, discovering truth through the heart rather than the intellect.",
            },
            10: {
                "Professional sensitivity may wound you, leaving you fearful of criticism or undervaluing your leadership. You may feel too vulnerable in public roles.",
                "Inverted, you lead through empathic strategy. You wield sensitivity as a strength, guiding others with compassion and emotional intelligence.",
            },
            11: {
                "Social circle emotional needs may wound you, leaving you fearful of rejection or drained by caretaking roles. You may feel invisible or unappreciated.",
                "Inverted, you create emotional sanctuaries in groups. You become the heart of communities, offering safety and belonging through your nurturing presence.",
            },
            12: {
                "Spiritual absorption may wound you, leaving you fearful of losing yourself in mystical experiences or overwhelmed by the unconscious. You may struggle with boundaries in spiritual practice.",
                "Inverted, you dissolve into the cosmic womb with strength. You embrace unity with the divine, wielding emotional surrender as a path to transcendence.",
            },
        },
        "Leo": {
            1: {
                "You may struggle with ego vulnerability, fearing rejection or feeling wounded when your self-expression isn’t validated. This can manifest as insecurity about your worth or constant need for recognition.",
                "Inverted, you embrace the self as a solar center. You radiate confidence and warmth, inspiring others through your authentic presence and fearless self-expression.",
            },
            2: {
                "Creative expression tied to validation may wound you, leaving you fearful of creating without praise. You may feel blocked when your art isn’t acknowledged.",
                "Inverted, you create because you must, not for applause. Your creativity becomes sacred ritual, shining regardless of external recognition.",
            },
            3: {
                "Dramatic communication may wound you, leaving you fearful of being dismissed as excessive or misunderstood. You may struggle with balancing passion and clarity.",
                "Inverted, every word becomes a performance of truth. You inspire others with your dramatic flair, turning communication into art that captivates and persuades.",
            },
            4: {
                "Family recognition needs may wound you, leaving you anxious about being seen or valued within your lineage. You may feel overshadowed or unappreciated.",
                "Inverted, you shine as your family’s brightest star. You embrace leadership within your lineage, offering courage and inspiration to those who came before and after.",
            },
            5: {
                "Creative performance anxiety may wound you, leaving you fearful of failure or hesitant to share your talents. You may feel blocked by perfectionism.",
                "Inverted, every act becomes sacred ritual. You transform performance into devotion, inspiring others through your courage to create authentically.",
            },
            6: {
                "Workplace pride issues may wound you, leaving you fearful of criticism or undervaluing your contributions. You may struggle with humility or recognition.",
                "Inverted, excellence becomes your natural state. You embody pride as strength, showing that confidence and mastery uplift everyone around you.",
            },
            7: {
                "Partnership ego clashes may wound you, leaving you fearful of losing individuality or struggling with dominance. You may feel torn between self-expression and compromise.",
                "Inverted, you find co-stars, not audiences. You thrive in partnerships that celebrate individuality, creating relationships where both shine equally.",
            },
            8: {
                "Transformation of pride may wound you, leaving you fearful of surrender or resistant to vulnerability. You may struggle with letting go of ego in deep change.",
                "Inverted, you die to be reborn more radiant. You embrace transformation as a path to greater brilliance, wielding pride as fuel for renewal.",
            },
            9: {
                "Belief system theatrics may wound you, leaving you fearful of being dismissed or misunderstood in spiritual or philosophical pursuits. You may feel pressured to perform faith.",
                "Inverted, faith becomes dramatic revelation. You inspire others by embodying belief with passion, turning spirituality into radiant expression.",
            },
            10: {
                "Career recognition obsession may wound you, leaving you fearful of invisibility or undervaluing your achievements. You may feel trapped by ambition for external validation.",
                "Inverted, you build legacy that outshines all. You embrace career as a stage for authentic brilliance, proving that true recognition comes from inner radiance.",
            },
            11: {
                "Social circle center stage needs may wound you, leaving you fearful of rejection or anxious about belonging. You may feel pressured to always perform.",
                "Inverted, you become the natural gravitational center. You inspire communities through your warmth and charisma, drawing people together effortlessly.",
            },
            12: {
                "Spiritual pride may wound you, leaving you fearful of surrender or resistant to humility in mystical practice. You may struggle with ego in spiritual growth.",
                "Inverted, union with the divine becomes ultimate performance. You embrace spirituality as radiant expression, embodying pride as devotion to cosmic truth.",
            },
        },
        "Virgo": {
            1: {
                "You may struggle with self‑criticism, constantly analyzing your flaws and feeling wounded by imperfection. This can manifest as anxiety about never being good enough.",
                "Inverted, you perfect the vessel of self. You transform critique into refinement, becoming a model of discipline and integrity that inspires others.",
            },
            2: {
                "Resource anxiety may wound you, leaving you fearful of scarcity or obsessed with managing every detail. You may feel burdened by the weight of responsibility.",
                "Inverted, you cultivate wealth through meticulous management. Your careful stewardship ensures abundance that is sustainable and secure.",
            },
            3: {
                "Communication perfectionism may wound you, leaving you fearful of speaking unless every word is flawless. You may feel silenced by your own standards.",
                "Inverted, your words become precision tools. You wield language with clarity and accuracy, offering insights that cut through confusion.",
            },
            4: {
                "Family duty burdens may wound you, leaving you overwhelmed by obligations or expectations. You may feel trapped in cycles of service.",
                "Inverted, you serve lineage through purification. You transform duty into devotion, becoming the healer who cleanses ancestral wounds.",
            },
            5: {
                "Creative inhibition may wound you, leaving you fearful of imperfection in art. You may hesitate to share your creations.",
                "Inverted, you craft as sacred geometry. Your creativity becomes precise and intentional, producing works that embody harmony and order.",
            },
            6: {
                "Workplace anxiety may wound you, leaving you fearful of mistakes or overwhelmed by details. You may feel undervalued despite your diligence.",
                "Inverted, you transform service into ritual. Your work becomes sacred practice, elevating routine into meaningful contribution.",
            },
            7: {
                "Relationship analysis may wound you, leaving you fearful of flaws or overly critical of partners. You may struggle to relax into intimacy.",
                "Inverted, you treat partnership as perfect system. You build relationships on clarity and mutual refinement, ensuring balance and growth.",
            },
            8: {
                "Transformation through analysis may wound you, leaving you fearful of surrender or trapped in overthinking. You may resist emotional depth.",
                "Inverted, you dissect to understand rebirth. You bring clarity to transformation, guiding others through change with wisdom and precision.",
            },
            9: {
                "Belief system skepticism may wound you, leaving you fearful of faith or dismissive of intuition. You may struggle to trust what cannot be proven.",
                "Inverted, you embrace faith through empirical evidence. You show that spirituality and reason can coexist, grounding belief in lived experience.",
            },
            10: {
                "Career perfectionism may wound you, leaving you fearful of failure or obsessed with flawless achievement. You may feel paralyzed by high standards.",
                "Inverted, you build flawless reputation. Your dedication to excellence becomes your strength, inspiring trust and respect in professional realms.",
            },
            11: {
                "Social improvement anxiety may wound you, leaving you fearful of imperfection in groups or overwhelmed by responsibility. You may feel burdened by collective flaws.",
                "Inverted, you perfect the collective. You become the reformer who elevates communities, ensuring growth through careful refinement.",
            },
            12: {
                "Spiritual materialism may wound you, leaving you fearful of disorder or clinging to ritual without deeper meaning. You may struggle with surrender.",
                "Inverted, you find sacredness in perfect order. You embody spirituality through discipline, showing that structure can be divine.",
            },
        },
        "Libra": {
            1: {
                "You may struggle with indecisive self‑image, feeling wounded by uncertainty about who you are. This can manifest as hesitation, fear of imbalance, or constant comparison with others.",
                "Inverted, you embrace balance as strategic advantage. You learn to wield harmony as power, showing that identity can be fluid yet strong.",
            },
            2: {
                "Value through relationships may wound you, leaving you fearful of being defined only by others. You may feel insecure when alone or undervalued outside of partnership.",
                "Inverted, you discover self‑worth independent of others. You cultivate inner balance, proving that relationships enhance rather than define your value.",
            },
            3: {
                "Diplomatic communication may wound you, leaving you fearful of conflict or silenced by the need to please. You may struggle to assert truth directly.",
                "Inverted, your words become peace treaties. You wield diplomacy as strength, creating dialogue that heals and unites.",
            },
            4: {
                "Family harmony obsession may wound you, leaving you anxious about conflict or burdened by the need to mediate. You may feel trapped by expectations of peacekeeping.",
                "Inverted, you balance lineage energies. You become the mediator who transforms family discord into harmony, ensuring growth through fairness.",
            },
            5: {
                "Creative partnership needs may wound you, leaving you fearful of creating alone or dependent on collaboration. You may feel blocked without external validation.",
                "Inverted, you create beauty through balance. Your art thrives in collaboration, weaving harmony into every expression.",
            },
            6: {
                "Workplace conflict avoidance may wound you, leaving you fearful of confrontation or undervaluing your contributions. You may struggle to assert yourself.",
                "Inverted, you master the art of equitable exchange. You transform workplaces into balanced ecosystems, ensuring fairness and cooperation.",
            },
            7: {
                "Relationship imbalance fears may wound you, leaving you fearful of unequal dynamics or anxious about dependency. You may feel trapped in cycles of compromise.",
                "Inverted, you master power dynamics. You build partnerships rooted in fairness, proving that equality is the foundation of love.",
            },
            8: {
                "Transformation through partnership may wound you, leaving you fearful of losing yourself in intimacy or resistant to shared change. You may struggle with vulnerability.",
                "Inverted, you embrace death and rebirth together. You wield transformation as shared strength, proving that vulnerability deepens connection.",
            },
            9: {
                "Belief system fairness may wound you, leaving you fearful of injustice or overwhelmed by contradictions. You may struggle to reconcile ideals with reality.",
                "Inverted, you embrace justice as divine principle. You become the seeker who finds truth through fairness, embodying balance in philosophy.",
            },
            10: {
                "Career diplomacy may wound you, leaving you fearful of conflict or undervaluing ambition. You may hesitate to assert leadership.",
                "Inverted, you succeed through strategic alliances. You wield diplomacy as strength, building careers on cooperation and fairness.",
            },
            11: {
                "Social harmony needs may wound you, leaving you fearful of rejection or burdened by the need to please. You may feel anxious about belonging.",
                "Inverted, you orchestrate perfect social symphony. You become the harmonizer who builds communities through fairness and balance.",
            },
            12: {
                "Spiritual balance may wound you, leaving you fearful of extremes or hesitant to surrender. You may struggle with integrating opposites.",
                "Inverted, you embrace equilibrium with the void. You wield balance as sacred practice, finding divinity in harmony itself.",
            },
        },
        "Scorpio": {
            1: {
                "You may struggle with intensity in self‑expression, feeling wounded by the depth of your emotions or the fear of being too much for others. This can manifest as secrecy or self‑protection.",
                "Inverted, you embrace power as your inherent nature. You radiate authenticity and strength, showing that intensity is a gift rather than a burden.",
            },
            2: {
                "Possessiveness over resources may wound you, leaving you fearful of loss or overly controlling of material security. You may feel anxious about scarcity or betrayal.",
                "Inverted, you cultivate wealth through strategic control. You learn to steward resources wisely, building abundance through discipline and foresight.",
            },
            3: {
                "Manipulative communication may wound you, leaving you fearful of being misunderstood or mistrusted. You may struggle with expressing truth directly.",
                "Inverted, your words penetrate to truth. You wield language with precision and depth, guiding others to clarity through honesty and insight.",
            },
            4: {
                "Family power dynamics may wound you, leaving you fearful of betrayal or burdened by secrets. You may feel trapped in cycles of control or manipulation.",
                "Inverted, you master lineage secrets. You transform ancestral wounds into wisdom, becoming the healer who brings hidden truths into light.",
            },
            5: {
                "Creative intensity may wound you, leaving you fearful of exposing your shadow or hesitant to share art that feels raw. You may struggle with vulnerability in expression.",
                "Inverted, you create from shadow depths. Your art becomes transformative, channeling intensity into works that inspire and heal.",
            },
            6: {
                "Workplace power struggles may wound you, leaving you fearful of betrayal or resistant to authority. You may feel trapped in cycles of conflict.",
                "Inverted, you control through understanding systems. You wield insight as strength, transforming workplaces by mastering hidden dynamics.",
            },
            7: {
                "Fear of betrayal in partnership may wound you, leaving you anxious about intimacy or resistant to vulnerability. You may struggle with trust.",
                "Inverted, you engineer betrayal into transformation. You rise in the ruin, proving that vulnerability can be rebirth and strength.",
            },
            8: {
                "Transformation obsession may wound you, leaving you fearful of surrender or overwhelmed by intensity. You may resist change even when it is necessary.",
                "Inverted, you master death to master life. You embrace transformation as sacred, wielding rebirth as a path to empowerment.",
            },
            9: {
                "Belief system intensity may wound you, leaving you fearful of surrender or resistant to faith. You may struggle with extremes in philosophy.",
                "Inverted, you embrace faith as total surrender. You embody devotion with depth, showing that intensity can be sacred truth.",
            },
            10: {
                "Career power ambitions may wound you, leaving you fearful of failure or obsessed with control. You may feel burdened by the need to dominate.",
                "Inverted, you build empire from ashes. You transform ambition into resilience, proving that true power comes from renewal.",
            },
            11: {
                "Social transformation may wound you, leaving you fearful of rejection or burdened by the need to control groups. You may feel anxious about belonging.",
                "Inverted, you remake circles in your image. You become the catalyst for collective rebirth, inspiring communities through transformation.",
            },
            12: {
                "Spiritual underworld navigation may wound you, leaving you fearful of hidden realms or resistant to surrender. You may struggle with mystical intensity.",
                "Inverted, you master all hidden realms. You embrace shadow as sacred, wielding spiritual depth as a path to transcendence.",
            },
        },
        "Sagittarius": {
            1: {
                "You may struggle with restless identity, feeling wounded by the need for constant expansion or fear of being confined. This can manifest as difficulty committing to one path or identity.",
                "Inverted, you embrace freedom as your essence. You show that identity can be vast and evolving, inspiring others to seek growth without fear.",
            },
            2: {
                "Financial recklessness may wound you, leaving you fearful of scarcity or guilty about indulgence. You may struggle with balancing adventure and stability.",
                "Inverted, you cultivate abundance through exploration. You discover wealth in experiences and opportunities, proving that prosperity comes from openness to the world.",
            },
            3: {
                "Overzealous communication may wound you, leaving you fearful of being dismissed as preachy or misunderstood. You may struggle with balancing passion and listening.",
                "Inverted, your words become arrows of truth. You inspire others with conviction, guiding them toward wisdom through your expansive vision.",
            },
            4: {
                "Family restlessness may wound you, leaving you fearful of confinement or disconnected from roots. You may feel torn between home and adventure.",
                "Inverted, you expand lineage horizons. You become the explorer who brings new wisdom to family, enriching roots with global perspective.",
            },
            5: {
                "Creative excess may wound you, leaving you fearful of being too much or scattered in expression. You may struggle with focus.",
                "Inverted, your art becomes a journey. You create expansively, inspiring others through works that embody freedom and exploration.",
            },
            6: {
                "Workplace restlessness may wound you, leaving you fearful of routine or resistant to structure. You may feel trapped in repetitive tasks.",
                "Inverted, you bring adventure into service. You transform work into exploration, inspiring innovation and growth through curiosity.",
            },
            7: {
                "Relationship freedom fears may wound you, leaving you anxious about commitment or resistant to intimacy. You may struggle with balancing independence and partnership.",
                "Inverted, you embrace partnership as shared adventure. You build relationships rooted in exploration, proving that love thrives on freedom.",
            },
            8: {
                "Transformation excess may wound you, leaving you fearful of intensity or overwhelmed by change. You may resist surrender to deep shifts.",
                "Inverted, you master rebirth as expansion. You embrace transformation as a journey, wielding change as a path to wisdom.",
            },
            9: {
                "Belief system dogmatism may wound you, leaving you fearful of questioning or rigid in philosophy. You may cling to certainty as safety.",
                "Inverted, you embrace truth as infinite horizon. You thrive in exploration of wisdom, showing that philosophy is a journey, not a destination.",
            },
            10: {
                "Career restlessness may wound you, leaving you fearful of confinement or undervaluing stability. You may struggle with long‑term goals.",
                "Inverted, you build legacy through exploration. Your career becomes a testament to freedom, proving that success can be expansive and adventurous.",
            },
            11: {
                "Social circle wanderlust may wound you, leaving you fearful of belonging or resistant to commitment in groups. You may feel disconnected.",
                "Inverted, you create global tribes. You build communities through exploration, connecting diverse people with shared vision and adventure.",
            },
            12: {
                "Spiritual excess may wound you, leaving you fearful of surrender or overwhelmed by mystical intensity. You may struggle with grounding.",
                "Inverted, you embrace the cosmos as sacred journey. You thrive spiritually by exploring infinite horizons, finding divinity in freedom itself.",
            },
        },
        "Capricorn": {
            1: {
                "You may struggle with self‑worth tied to achievement, feeling wounded when progress is slow or recognition is withheld. This can manifest as insecurity about your identity without success.",
                "Inverted, you embody discipline as identity. You show that persistence and responsibility are strengths, inspiring others through your steady presence.",
            },
            2: {
                "Resource anxiety may wound you, leaving you fearful of scarcity or overly focused on material accumulation. You may feel burdened by the weight of responsibility.",
                "Inverted, you cultivate wealth through patience and structure. You build resources steadily, ensuring abundance that endures across time.",
            },
            3: {
                "Rigid communication may wound you, leaving you fearful of speaking unless certain or authoritative. You may struggle with flexibility in dialogue.",
                "Inverted, your words carry weight and authority. You speak with clarity and discipline, offering wisdom that others respect and rely upon.",
            },
            4: {
                "Family duty burdens may wound you, leaving you overwhelmed by obligations or expectations. You may feel trapped in cycles of responsibility.",
                "Inverted, you build legacy through devotion. You transform duty into sacred service, becoming the pillar of strength for your lineage.",
            },
            5: {
                "Creative inhibition may wound you, leaving you fearful of imperfection or hesitant to share your art. You may feel blocked by high standards.",
                "Inverted, you create with discipline and mastery. Your art becomes timeless, embodying structure and endurance that inspire generations.",
            },
            6: {
                "Workplace rigidity may wound you, leaving you fearful of change or resistant to innovation. You may feel trapped in repetitive tasks.",
                "Inverted, you master systems through discipline. You transform work into sacred practice, ensuring stability and growth through persistence.",
            },
            7: {
                "Relationship duty may wound you, leaving you fearful of imbalance or burdened by responsibility. You may struggle with intimacy when obligations dominate.",
                "Inverted, you build partnerships on loyalty and endurance. You prove that commitment and responsibility are foundations of lasting love.",
            },
            8: {
                "Transformation resistance may wound you, leaving you fearful of surrender or reluctant to let go. You may resist deep change.",
                "Inverted, you embrace transformation as structured renewal. You rebuild steadily, proving that rebirth can be disciplined and enduring.",
            },
            9: {
                "Belief system rigidity may wound you, leaving you fearful of questioning or resistant to new philosophies. You may cling to tradition as safety.",
                "Inverted, you embody wisdom through structure. You show that philosophy can be disciplined, grounding truth in lived experience and responsibility.",
            },
            10: {
                "Career ambition obsession may wound you, leaving you fearful of failure or burdened by responsibility. You may feel trapped by external expectations.",
                "Inverted, you build legacy through persistence. You transform ambition into resilience, proving that true success comes from endurance and discipline.",
            },
            11: {
                "Social duty burdens may wound you, leaving you fearful of rejection or overwhelmed by responsibility in groups. You may feel anxious about belonging.",
                "Inverted, you anchor communities through responsibility. You become the pillar of collective strength, ensuring stability and growth for all.",
            },
            12: {
                "Spiritual rigidity may wound you, leaving you fearful of surrender or resistant to mystical experiences. You may struggle with flexibility in spiritual practice.",
                "Inverted, you embody discipline as sacred devotion. You show that structure can be divine, grounding spirituality in persistence and endurance.",
            },
        },
        "Aquarius": {
            1: {
                "You may struggle with feeling alienated or misunderstood, wounded by the sense that your individuality sets you apart from others. This can manifest as loneliness or fear of rejection.",
                "Inverted, you embrace uniqueness as your greatest gift. You show others that individuality is strength, inspiring communities through your authenticity.",
            },
            2: {
                "Resource detachment may wound you, leaving you fearful of scarcity or disconnected from material needs. You may undervalue stability in pursuit of ideals.",
                "Inverted, you cultivate wealth through innovation. You harness unconventional methods to build abundance, proving that creativity can generate security.",
            },
            3: {
                "Unconventional communication may wound you, leaving you fearful of being misunderstood or dismissed. You may struggle with expressing radical ideas.",
                "Inverted, your words become sparks of revolution. You inspire others with visionary language, turning unconventional thought into collective progress.",
            },
            4: {
                "Family alienation may wound you, leaving you fearful of rejection or disconnected from roots. You may feel burdened by being different within your lineage.",
                "Inverted, you expand family horizons. You bring innovation and new perspectives to your lineage, transforming tradition into evolution.",
            },
            5: {
                "Creative eccentricity may wound you, leaving you fearful of ridicule or hesitant to share unconventional art. You may feel blocked by self‑doubt.",
                "Inverted, you create from radical vision. Your art becomes revolutionary, inspiring others through originality and boldness.",
            },
            6: {
                "Workplace nonconformity may wound you, leaving you fearful of rejection or undervalued for your unconventional methods. You may struggle with rigid systems.",
                "Inverted, you transform workplaces through innovation. You show that progress comes from breaking norms, inspiring change through visionary ideas.",
            },
            7: {
                "Relationship detachment may wound you, leaving you fearful of intimacy or resistant to vulnerability. You may struggle with balancing independence and connection.",
                "Inverted, you embrace partnership as shared evolution. You build relationships rooted in freedom and growth, proving that love thrives on individuality.",
            },
            8: {
                "Transformation through detachment may wound you, leaving you fearful of surrender or resistant to emotional depth. You may struggle with vulnerability in change.",
                "Inverted, you master transformation through innovation. You embrace rebirth as evolution, wielding detachment as clarity in deep change.",
            },
            9: {
                "Belief system radicalism may wound you, leaving you fearful of rejection or resistant to tradition. You may struggle with integrating unconventional philosophies.",
                "Inverted, you embody wisdom through innovation. You show that philosophy evolves through radical thought, inspiring others with visionary beliefs.",
            },
            10: {
                "Career unconventionality may wound you, leaving you fearful of rejection or undervalued for radical ambition. You may struggle with recognition in traditional systems.",
                "Inverted, you succeed through innovation. You build careers on visionary ideas, proving that unconventional paths lead to progress.",
            },
            11: {
                "Social alienation may wound you, leaving you fearful of rejection or disconnected from groups. You may feel misunderstood in communities.",
                "Inverted, you build visionary tribes. You create communities through innovation, inspiring collective progress and unity.",
            },
            12: {
                "Spiritual detachment may wound you, leaving you fearful of surrender or disconnected from mystical experiences. You may struggle with grounding in spiritual practice.",
                "Inverted, you embrace cosmic innovation. You discover divinity through radical thought, embodying spirituality as visionary evolution.",
            },
        },
        "Pisces": {
            1: {
                "You may struggle with dissolving identity, feeling wounded by uncertainty about boundaries or fear of losing yourself. This can manifest as confusion or vulnerability in self‑expression.",
                "Inverted, you embrace unity as identity. You show that self can be infinite, inspiring others through compassion and spiritual presence.",
            },
            2: {
                "Resource confusion may wound you, leaving you fearful of scarcity or disconnected from material needs. You may struggle with grounding abundance.",
                "Inverted, you cultivate wealth through surrender. You discover prosperity in flow, proving that abundance comes from trust in the universe.",
            },
            3: {
                "Communication vagueness may wound you, leaving you fearful of being misunderstood or dismissed. You may struggle with clarity in dialogue.",
                "Inverted, your words become poetry of the soul. You inspire others with mystical language, turning communication into art that transcends logic.",
            },
            4: {
                "Family dissolution may wound you, leaving you fearful of instability or disconnected from roots. You may feel burdened by confusion in lineage.",
                "Inverted, you embrace family as spiritual union. You transform lineage wounds into compassion, creating homes that embody unconditional love.",
            },
            5: {
                "Creative confusion may wound you, leaving you fearful of imperfection or hesitant to share art. You may feel blocked by lack of clarity.",
                "Inverted, you create from mystical flow. Your art becomes transcendent, inspiring others through imagination and spiritual depth.",
            },
            6: {
                "Workplace vagueness may wound you, leaving you fearful of being undervalued or misunderstood. You may struggle with structure in service.",
                "Inverted, you transform service into compassion. You embody empathy in work, showing that care and intuition are strengths.",
            },
            7: {
                "Relationship dissolution may wound you, leaving you fearful of abandonment or resistant to intimacy. You may struggle with boundaries in love.",
                "Inverted, you embrace love as spiritual union. You build partnerships rooted in compassion, proving that vulnerability is sacred strength.",
            },
            8: {
                "Transformation confusion may wound you, leaving you fearful of surrender or overwhelmed by intensity. You may resist deep change.",
                "Inverted, you master rebirth through surrender. You embrace transformation as mystical renewal, wielding compassion as power.",
            },
            9: {
                "Belief system vagueness may wound you, leaving you fearful of faith or resistant to clarity. You may struggle with grounding philosophy.",
                "Inverted, you embody wisdom through mystical truth. You show that faith can be infinite, inspiring others through compassion and imagination.",
            },
            10: {
                "Career dissolution may wound you, leaving you fearful of instability or undervaluing ambition. You may struggle with recognition in public life.",
                "Inverted, you succeed through compassion and vision. You build careers on empathy, proving that success can be mystical and transcendent.",
            },
            11: {
                "Social confusion may wound you, leaving you fearful of rejection or disconnected from groups. You may feel invisible or misunderstood.",
                "Inverted, you create communities through compassion. You inspire collective unity, building tribes that thrive on empathy and imagination.",
            },
            12: {
                "Spiritual overwhelm may wound you, leaving you fearful of surrender or resistant to mystical intensity. You may struggle with boundaries in spiritual practice.",
                "Inverted, you dissolve into cosmic unity. You embrace spirituality as infinite compassion, embodying transcendence through surrender and love.",
            },
        },
    } // end of interpretations map

    // Default fallback if sign/house not found
    if houses, ok := interpretations[sign]; ok {
        if pair, ok := houses[house]; ok {
            return pair[0], pair[1]
        }
    }
    return "No interpretation available.", "No strength available."
}
