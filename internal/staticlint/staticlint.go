// staticlint - пакет с кастомной реализацией multichecker, а также анализатора NoDirectOsExitAnalyzer.
package staticlint

import (
	"github.com/kisielk/errcheck/errcheck"
	"github.com/timakin/bodyclose/passes/bodyclose"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/asmdecl"
	"golang.org/x/tools/go/analysis/passes/assign"
	"golang.org/x/tools/go/analysis/passes/atomic"
	"golang.org/x/tools/go/analysis/passes/bools"
	"golang.org/x/tools/go/analysis/passes/buildtag"
	"golang.org/x/tools/go/analysis/passes/cgocall"
	"golang.org/x/tools/go/analysis/passes/composite"
	"golang.org/x/tools/go/analysis/passes/copylock"
	"golang.org/x/tools/go/analysis/passes/ctrlflow"
	"golang.org/x/tools/go/analysis/passes/deepequalerrors"
	"golang.org/x/tools/go/analysis/passes/errorsas"
	"golang.org/x/tools/go/analysis/passes/fieldalignment"
	"golang.org/x/tools/go/analysis/passes/findcall"
	"golang.org/x/tools/go/analysis/passes/framepointer"
	"golang.org/x/tools/go/analysis/passes/httpresponse"
	"golang.org/x/tools/go/analysis/passes/ifaceassert"
	"golang.org/x/tools/go/analysis/passes/loopclosure"
	"golang.org/x/tools/go/analysis/passes/lostcancel"
	"golang.org/x/tools/go/analysis/passes/nilfunc"
	"golang.org/x/tools/go/analysis/passes/nilness"
	"golang.org/x/tools/go/analysis/passes/pkgfact"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/reflectvaluecompare"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/shift"
	"golang.org/x/tools/go/analysis/passes/sigchanyzer"
	"golang.org/x/tools/go/analysis/passes/sortslice"
	"golang.org/x/tools/go/analysis/passes/stdmethods"
	"golang.org/x/tools/go/analysis/passes/stringintconv"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"golang.org/x/tools/go/analysis/passes/testinggoroutine"
	"golang.org/x/tools/go/analysis/passes/tests"
	"golang.org/x/tools/go/analysis/passes/unmarshal"
	"golang.org/x/tools/go/analysis/passes/unreachable"
	"golang.org/x/tools/go/analysis/passes/unsafeptr"
	"golang.org/x/tools/go/analysis/passes/unusedresult"
	"golang.org/x/tools/go/analysis/passes/unusedwrite"
	"honnef.co/go/tools/staticcheck"
)

// RunLint запускает multichecker c набором анализаторов.
func RunLint() {
	var analyzers = []*analysis.Analyzer{
		// NoDirectOsExitAnalyzer проверяет, что в функции main пакета main
		// нет прямых вызовов os.Exit (самописный анализатор)
		NoDirectOsExitAnalyzer,

		// bodyclose проверяет, закрывается ли response.Body после HTTP-запросов.
		// (публичный анализатор)
		bodyclose.Analyzer,

		// errcheck проверяет необработанные ошибки
		// (публичный анализатор)
		errcheck.Analyzer,

		// asmdecl проверяет корректность объявлений Go-ассемблера.
		asmdecl.Analyzer,

		// assign обнаруживает бесполезные присваивания (например, x = x).
		assign.Analyzer,

		// atomic проверяет некорректное использование sync/atomic.
		atomic.Analyzer,

		// bools выявляет подозрительные операции с bool (например, x == true).
		bools.Analyzer,

		// buildtag проверяет корректность //go:build тегов.
		buildtag.Analyzer,

		// cgocall обнаруживает нарушения правил вызовов CGO.
		cgocall.Analyzer,

		// composite проверяет неявное использование композитных литералов.
		composite.Analyzer,

		// copylock обнаруживает копирование мьютексов и других lock-объектов.
		copylock.Analyzer,

		// ctrlflow анализирует поток управления (используется другими анализаторами).
		ctrlflow.Analyzer,

		// deepequalerrors выявляет использование reflect.DeepEqual с ошибками.
		deepequalerrors.Analyzer,

		// errorsas проверяет корректность использования errors.As.
		errorsas.Analyzer,

		// fieldalignment предлагает оптимизацию размера структур.
		fieldalignment.Analyzer,

		// findcall ищет вызовы функций с определенными сигнатурами.
		findcall.Analyzer,

		// framepointer анализирует использование указателей на фреймы.
		framepointer.Analyzer,

		// httpresponse проверяет обработку HTTP-ответов.
		httpresponse.Analyzer,

		// ifaceassert обнаруживает некорректные утверждения типов (type assertion).
		ifaceassert.Analyzer,

		// loopclosure выявляет проблемы с замыканиями в циклах.
		loopclosure.Analyzer,

		// lostcancel проверяет утечку контекстов (невызов cancel()).
		lostcancel.Analyzer,

		// nilfunc обнаруживает сравнение функций с nil.
		nilfunc.Analyzer,

		// nilness анализирует возможные nil-паники.
		nilness.Analyzer,

		// pkgfact собирает факты о пакетах (для сложных анализаторов).
		pkgfact.Analyzer,

		// printf проверяет формат-строки в Printf-подобных функциях.
		printf.Analyzer,

		// reflectvaluecompare выявляет сравнение reflect.Value с ==.
		reflectvaluecompare.Analyzer,

		// shadow обнаруживает перекрытие переменных (variable shadowing).
		shadow.Analyzer,

		// shift проверяет некорректные операции битового сдвига.
		shift.Analyzer,

		// sigchanyzer выявляет неправильное использование chan os.Signal.
		sigchanyzer.Analyzer,

		// sortslice проверяет корректность реализации sort.Interface.
		sortslice.Analyzer,

		// stdmethods ищет ошибки в методах, удовлетворяющих стандартным интерфейсам.
		stdmethods.Analyzer,

		// stringintconv обнаруживает неявное преобразование string <-> int.
		stringintconv.Analyzer,

		// structtag проверяет синтаксис тегов структур.
		structtag.Analyzer,

		// testinggoroutine выявляет утечку горутин в тестах.
		testinggoroutine.Analyzer,

		// tests обнаруживает распространенные ошибки в тестах.
		tests.Analyzer,

		// unmarshal проверяет корректность использования unmarshal-функций.
		unmarshal.Analyzer,

		// unreachable выявляет недостижимый код.
		unreachable.Analyzer,

		// unsafeptr обнаруживает некорректное использование unsafe.Pointer.
		unsafeptr.Analyzer,

		// unusedresult проверяет игнорирование возвращаемых ошибок.
		unusedresult.Analyzer,

		// unusedwrite выявляет запись в переменные без последующего чтения.
		unusedwrite.Analyzer,
	}

	for _, analyzer := range staticcheck.Analyzers {
		// Добавляем все SA-анализаторы из пакета staticcheck (их значение можно найти в документации https://staticcheck.dev/docs/checks/).
		if analyzer.Analyzer.Name[:2] == "SA" ||
			// Также добавляем анализатор ST1005, который ищет некорректно отформатированные строки ошибок.
			analyzer.Analyzer.Name == "ST1005" {
			analyzers = append(analyzers, analyzer.Analyzer)
		}
	}

	multichecker.Main(
		analyzers...,
	)
}
